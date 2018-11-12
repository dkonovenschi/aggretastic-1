package aggretastic

import "github.com/olivere/elastic"

// StatsAggregation is a multi-value metrics aggregation that computes stats
// over numeric values extracted from the aggregated documents.
// These values can be extracted either from specific numeric fields
// in the documents, or be generated by a provided script.
// See: https://www.elastic.co/guide/en/elasticsearch/reference/6.2/search-aggregations-metrics-stats-aggregation.html
type StatsAggregation struct {
	*aggregation

	field  string
	script *elastic.Script
	format string
	meta   map[string]interface{}
}

func NewStatsAggregation() *StatsAggregation {
	a := &StatsAggregation{}
	a.aggregation = nilAggregation()

	return a
}

func (a *StatsAggregation) Field(field string) *StatsAggregation {
	a.field = field
	return a
}

func (a *StatsAggregation) Script(script *elastic.Script) *StatsAggregation {
	a.script = script
	return a
}

func (a *StatsAggregation) Format(format string) *StatsAggregation {
	a.format = format
	return a
}

func (a *StatsAggregation) SubAggregation(name string, subAggregation Aggregation) *StatsAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *StatsAggregation) Meta(metaData map[string]interface{}) *StatsAggregation {
	a.meta = metaData
	return a
}

func (a *StatsAggregation) Source() (interface{}, error) {
	// Example:
	//	{
	//    "aggs" : {
	//      "grades_stats" : { "stats" : { "field" : "grade" } }
	//    }
	//	}
	// This method returns only the { "stats" : { "field" : "grade" } } part.

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["stats"] = opts

	// ValuesSourceAggregationBuilder
	if a.field != "" {
		opts["field"] = a.field
	}
	if a.script != nil {
		src, err := a.script.Source()
		if err != nil {
			return nil, err
		}
		opts["script"] = src
	}
	if a.format != "" {
		opts["format"] = a.format
	}

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