package aggretastic

// SumBucketAggregation is a sibling pipeline aggregation which calculates
// the sum across all buckets of a specified metric in a sibling aggregation.
// The specified metric must be numeric and the sibling aggregation must
// be a multi-bucket aggregation.
//
// For more details, see
// https://www.elastic.co/guide/en/elasticsearch/reference/6.2/search-aggregations-pipeline-sum-bucket-aggregation.html
type SumBucketAggregation struct {
	*finiteAggregation

	format    string
	gapPolicy string

	meta         map[string]interface{}
	bucketsPaths []string
}

// NewSumBucketAggregation creates and initializes a new SumBucketAggregation.
func NewSumBucketAggregation() *SumBucketAggregation {
	a := &SumBucketAggregation{
		bucketsPaths: make([]string, 0),
	}
	a.finiteAggregation = newFiniteAggregation()

	return a
}

// Format to use on the output of this aggregation.
func (a *SumBucketAggregation) Format(format string) *SumBucketAggregation {
	a.format = format
	return a
}

// GapPolicy defines what should be done when a gap in the series is discovered.
// Valid values include "insert_zeros" or "skip". Default is "insert_zeros".
func (a *SumBucketAggregation) GapPolicy(gapPolicy string) *SumBucketAggregation {
	a.gapPolicy = gapPolicy
	return a
}

// GapInsertZeros inserts zeros for gaps in the series.
func (a *SumBucketAggregation) GapInsertZeros() *SumBucketAggregation {
	a.gapPolicy = "insert_zeros"
	return a
}

// GapSkip skips gaps in the series.
func (a *SumBucketAggregation) GapSkip() *SumBucketAggregation {
	a.gapPolicy = "skip"
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *SumBucketAggregation) Meta(metaData map[string]interface{}) *SumBucketAggregation {
	a.meta = metaData
	return a
}

// BucketsPath sets the paths to the buckets to use for this pipeline aggregator.
func (a *SumBucketAggregation) BucketsPath(bucketsPaths ...string) *SumBucketAggregation {
	a.bucketsPaths = append(a.bucketsPaths, bucketsPaths...)
	return a
}

// Source returns the a JSON-serializable interface.
func (a *SumBucketAggregation) Source() (interface{}, error) {
	source := make(map[string]interface{})
	params := make(map[string]interface{})
	source["sum_bucket"] = params

	if a.format != "" {
		params["format"] = a.format
	}
	if a.gapPolicy != "" {
		params["gap_policy"] = a.gapPolicy
	}

	// Add buckets paths
	switch len(a.bucketsPaths) {
	case 0:
	case 1:
		params["buckets_path"] = a.bucketsPaths[0]
	default:
		params["buckets_path"] = a.bucketsPaths
	}

	// Add Meta data if available
	if len(a.meta) > 0 {
		source["meta"] = a.meta
	}

	return source, nil
}