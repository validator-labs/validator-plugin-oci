package types

// UsageResult describes the maximum usage for an arbitrary category
type UsageResult struct {
	Description string
	MaxUsage    float64
}

// UsageMap maps categories to their usage
type UsageMap map[string]float64

// Max returns a UsageResult describing the category with the maximum usage within a UsageMap
func (u UsageMap) Max() *UsageResult {
	var maxUsage float64
	var maxUsageKey string
	for k, v := range u {
		if v > maxUsage {
			maxUsage = v
			maxUsageKey = k
		}
	}
	return &UsageResult{Description: maxUsageKey, MaxUsage: maxUsage}
}
