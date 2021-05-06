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

package explain

import (
	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/endink/go-sharding/core"
)

type brackets int

const (
	bracketsStart brackets = iota
	bracketsEnd
)

//logic stack, priority support
type LogicPriorityStack struct {
	valueStack *arraystack.Stack
	logicStack *arraystack.Stack
}

type valueAndProcessor struct {
	value     interface{}
	processor func(logic core.BinaryLogic, value interface{}) error
}

var placeholder = &valueAndProcessor{}

func newLogicPriorityStack() *LogicPriorityStack {
	return &LogicPriorityStack{
		valueStack: arraystack.New(),
		logicStack: arraystack.New(),
	}
}

func (ls *LogicPriorityStack) PushValue(v interface{}, processor func(logic core.BinaryLogic, value interface{}) error) error {
	if ls.current().logicStack.Size() > 0 {
		l, ok := ls.current().logicStack.Pop()
		if ok {
			logic := l.(core.BinaryLogic)
			switch logic {
			case core.LogicAnd:
				if pre, hasValue := ls.current().valueStack.Pop(); hasValue {
					vp := pre.(*valueAndProcessor)
					if vp != placeholder {
						if err := vp.processor(core.LogicAnd, vp.value); err != nil {
							return err
						}
					}
					if err := processor(core.LogicAnd, v); err != nil {
						return err
					}
					ls.current().valueStack.Push(placeholder) //计算结果也要入栈
				}
			case core.LogicOr:
				ls.current().valueStack.Push(&valueAndProcessor{
					value:     v,
					processor: processor,
				})
			}
		}
	} else {
		ls.current().valueStack.Push(&valueAndProcessor{
			value:     v,
			processor: processor,
		})
	}
	return nil
}

func (ls *LogicPriorityStack) current() *LogicPriorityStack {
	if v, ok := ls.valueStack.Peek(); ok {
		if lgStack, isLgStack := v.(*LogicPriorityStack); isLgStack {
			return lgStack
		}
	}
	return ls
}

func (ls *LogicPriorityStack) PushLogic(logic core.BinaryLogic) {
	ls.current().logicStack.Push(logic)
}

func (ls *LogicPriorityStack) PushBracketsStart() {
	ls.current().valueStack.Push(newLogicPriorityStack())
}

func (ls *LogicPriorityStack) PushBracketsEnd() error {
	if v, ok := ls.valueStack.Peek(); ok {
		if lgStack, isLgStack := v.(*LogicPriorityStack); isLgStack {
			ls.valueStack.Pop()
			return lgStack.Calc()
		}
	}
	return nil
}

func (ls *LogicPriorityStack) Calc() error {
	for !ls.valueStack.Empty() {
		v, _ := ls.valueStack.Pop()
		vp := v.(*valueAndProcessor)
		if vp == placeholder {
			continue
		}
		if err := vp.processor(core.LogicOr, vp.value); err != nil {
			return err
		}
	}
	return nil
}