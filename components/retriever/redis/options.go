/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package redis

import (
	"github.com/cloudwego/eino/components/retriever"
)

type implOptions struct {
	FilterQuery string
}

// WithFilterQuery redis filter query.
// see: https://redis.io/docs/latest/develop/interact/search-and-query/advanced-concepts/vectors/#filters
func WithFilterQuery(filter string) retriever.Option {
	return retriever.WrapImplSpecificOptFn(func(o *implOptions) {
		o.FilterQuery = filter
	})
}
