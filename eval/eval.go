package eval

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/antonmedv/expr"
	"github.com/pkg/errors"

	"github.com/runabol/tork/task"
)

var exprMatcher = regexp.MustCompile(`{{\s*(.+?)\s*}}`)

func EvaluateTask(t *task.Task, c map[string]any) error {
	// evaluate name
	name, err := evaluateTemplate(t.Name, c)
	if err != nil {
		return err
	}
	t.Name = name
	// evaluate var
	var_, err := evaluateTemplate(t.Var, c)
	if err != nil {
		return err
	}
	t.Var = var_
	// evaluate image
	img, err := evaluateTemplate(t.Image, c)
	if err != nil {
		return err
	}
	t.Image = img
	// evaluate queue
	q, err := evaluateTemplate(t.Queue, c)
	if err != nil {
		return err
	}
	t.Queue = q
	// evaluate the env vars
	env := t.Env
	for k, v := range env {
		result, err := evaluateTemplate(v, c)
		if err != nil {
			return err
		}
		env[k] = result
	}
	t.Env = env
	// evaluate if expr
	ifExpr, err := evaluateTemplate(t.If, c)
	if err != nil {
		return err
	}
	t.If = ifExpr
	// evaluate pre-tasks
	pres := make([]*task.Task, len(t.Pre))
	for i, pre := range t.Pre {
		if err := EvaluateTask(pre, c); err != nil {
			return err
		}
		pres[i] = pre
	}
	t.Pre = pres
	// evaluate post-tasks
	posts := make([]*task.Task, len(t.Post))
	for i, post := range t.Post {
		if err := EvaluateTask(post, c); err != nil {
			return err
		}
		posts[i] = post
	}
	t.Post = posts
	// evaluate parallel tasks
	parallel := make([]*task.Task, len(t.Parallel))
	for i, par := range t.Parallel {
		if err := EvaluateTask(par, c); err != nil {
			return err
		}
		parallel[i] = par
	}
	t.Parallel = parallel
	return nil
}

func evaluateTemplate(ex string, c map[string]any) (string, error) {
	if ex == "" {
		return "", nil
	}
	loc := 0
	var buf bytes.Buffer
	for _, match := range exprMatcher.FindAllStringSubmatchIndex(ex, -1) {
		startTag := match[0]
		endTag := match[1]
		startExpr := match[2]
		endExpr := match[3]
		buf.WriteString(ex[loc:startTag])
		ev, err := EvaluateExpr(ex[startExpr:endExpr], c)
		if err != nil {
			return "", err
		}
		buf.WriteString(fmt.Sprintf("%v", ev))
		loc = endTag
	}
	buf.WriteString(ex[loc:])
	return buf.String(), nil
}

func EvaluateExpr(ex string, c map[string]any) (any, error) {
	// if the expression is marked with {{ }}
	// we want to strip these off
	if matches := exprMatcher.FindStringSubmatch(ex); matches != nil {
		ex = matches[1]
	}
	env := map[string]any{
		"randomInt": randomInt,
		"coinflip":  coinflip,
		"range":     range_,
		"parseJSON": parseJSON,
	}
	for k, v := range c {
		env[k] = v
	}
	program, err := expr.Compile(ex, expr.Env(env))
	if err != nil {
		return "", errors.Wrapf(err, "error compiling expression: %s", ex)
	}
	output, err := expr.Run(program, env)
	if err != nil {
		return "", errors.Wrapf(err, "error evaluating expression: %s", ex)
	}
	return output, nil
}
