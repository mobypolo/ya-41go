package agent

var sources []MetricSource

func Register(source MetricSource) {
	sources = append(sources, source)
}

func CollectAll() ([]Metric, error) {
	var result []Metric
	for _, s := range sources {
		val, err := s.Collect()
		if err != nil {
			continue
		}
		result = append(result, Metric{
			Name:  s.Name(),
			Type:  s.Type(),
			Value: val,
		})
	}
	return result, nil
}
