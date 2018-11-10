package aggretastic

import "github.com/olivere/elastic"

// FilterAggregation defines a single bucket of all the documents
// in the current document set context that match a specified filter.
// Often this will be used to narrow down the current aggregation context
// to a specific set of documents.
// See: https://www.elastic.co/guide/en/elasticsearch/reference/6.2/search-aggregations-bucket-filter-aggregation.html
type FilterAggregation struct {
	*tree

	filter elastic.Query
	meta   map[string]interface{}
}

func NewFilterAggregation() *FilterAggregation {
	a := &FilterAggregation{}
	a.tree = nilAggregationTree(a)

	return a
}

func (a *FilterAggregation) SubAggregation(name string, subAggregation Aggregation) *FilterAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *FilterAggregation) Meta(metaData map[string]interface{}) *FilterAggregation {
	a.meta = metaData
	return a
}

func (a *FilterAggregation) Filter(filter elastic.Query) *FilterAggregation {
	a.filter = filter
	return a
}

func (a *FilterAggregation) Source() (interface{}, error) {
	// Example:
	//	{
	//    "aggs" : {
	//         "in_stock_products" : {
	//             "filter" : { "range" : { "stock" : { "gt" : 0 } } }
	//         }
	//    }
	//	}
	// This method returns only the { "filter" : {} } part.

	src, err := a.filter.Source()
	if err != nil {
		return nil, err
	}
	source := make(map[string]interface{})
	source["filter"] = src

	// AggregationBuilder (SubAggregations)
	if len(a.subAggregations) > 0 {
		aggsMap := make(map[string]interface{})
		source["aggregations"] = aggsMap
		for name, aggregate := range a.subAggregations {
			src, err := aggregate.Source()
			if err != nil {
				return nil, err
			}
			aggsMap[name] = src
		}
	}

	// Add Meta data if available
	if len(a.meta) > 0 {
		source["meta"] = a.meta
	}

	return source, nil
}
