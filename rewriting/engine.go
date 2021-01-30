/*
 * Copyright 2021. Go-Sharding Author All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 *  File author: Anders Xiao
 */

package rewriting

import (
	"github.com/XiaoMi/Gaea/explain"
	"github.com/pingcap/errors"
	"github.com/pingcap/parser/ast"
)

var _ explain.Rewriter = &Engine{}

type Engine struct {
	runtime Runtime
}

func NewEngine(runtime Runtime) *Engine {
	return &Engine{
		runtime: runtime,
	}
}

func (engine *Engine) RewriteTable(tableName *ast.TableName, explainContext explain.Context) (explain.RewriteNodeResult, error) {
	sd, ok, fe := explain.FindShardingTableByTable(tableName, explainContext, "")
	if fe != nil {
		return nil, fe
	}
	if ok {
		if writer, err := NewTableNameWriter(tableName, engine.runtime); err == nil {
			return explain.ResultFromNode(writer, sd.Name), nil
		} else {
			return nil, err
		}
	}
	return explain.NoneRewriteNodeResult, nil
}

func (engine *Engine) RewriteField(columnName *ast.ColumnNameExpr, explainContext explain.Context) (explain.RewriteExprResult, error) {
	sd, ok, err := explain.FindShardingTableByColumn(columnName, explainContext, false)
	if err != nil {
		return nil, err
	}
	if ok {
		if writer, e := NewColumnNameWriter(columnName, explainContext, engine.runtime, sd.Name); e == nil {
			return explain.ResultFromExrp(writer, sd.Name, explain.GetColumn(columnName.Name)), nil
		} else {
			return nil, e
		}
	}
	return explain.NoneRewriteExprResult, nil
}

func (engine *Engine) RewriteColumn(columnName *ast.ColumnNameExpr, explainContext explain.Context) (explain.RewriteExprResult, error) {
	sd, ok, err := explain.FindShardingTableByColumn(columnName, explainContext, true)
	if err != nil {
		return nil, err
	}
	if ok {
		if writer, e := NewColumnNameWriter(columnName, explainContext, engine.runtime, sd.Name); e == nil {
			return explain.ResultFromExrp(writer, sd.Name, explain.GetColumn(columnName.Name)), nil
		} else {
			return nil, e
		}
	}
	return explain.NoneRewriteExprResult, nil
}

func (engine *Engine) RewritePatterIn(patternIn *ast.PatternInExpr, explainContext explain.Context) (explain.RewriteExprResult, error) {
	columnNameExpr, ok := patternIn.Expr.(*ast.ColumnNameExpr)
	if !ok {
		return nil, errors.New("pattern in statement required ColumnNameExpr")
	}
	sd, ok, err := explain.FindShardingTableByColumn(columnNameExpr, explainContext, true)
	if err != nil {
		return nil, err
	}
	if ok {
		if writer, e := NewPatternInWriter(patternIn, explainContext, engine.runtime, sd); e == nil {
			return explain.ResultFromExrp(writer, sd.Name, explain.GetColumn(columnNameExpr.Name)), nil
		} else {
			return nil, e
		}
	}
	return explain.NoneRewriteExprResult, nil
}

func (engine *Engine) RewriteBetween(between *ast.BetweenExpr, explainContext explain.Context) (explain.RewriteExprResult, error) {
	columnNameExpr, ok := between.Expr.(*ast.ColumnNameExpr)
	if !ok {
		return nil, errors.New("between and statement required ColumnNameExpr")
	}
	sd, ok, err := explain.FindShardingTableByColumn(columnNameExpr, explainContext, true)
	if err != nil {
		return nil, err
	}
	if ok {
		if writer, e := NewBetweenWriter(between, explainContext, engine.runtime, sd); e == nil {
			return explain.ResultFromExrp(writer, sd.Name, explain.GetColumn(columnNameExpr.Name)), nil
		} else {
			return nil, e
		}
	}
	return explain.NoneRewriteExprResult, nil
}

func (engine *Engine) RewriteLimit(limit *ast.Limit, explainContext explain.Context) (explain.RewriteLimitResult, error) {
	lookup := explainContext.LimitLookup()
	if lookup.HasLimit() && lookup.HasOffset() {
		writer, err := NewLimitWriter(explainContext)
		if err != nil {
			return nil, err
		}
		return explain.ResultFromLimit(writer), nil
	}
	return explain.NoneRewriteLimitResult, nil
}