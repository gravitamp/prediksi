package main

import (
	"sort"
)

type sorted struct {
	category    string
	probability float64
}

type Classifier struct {
	datatrain           []Condition
	datatest            []Condition
	waktu               (map[string]map[string]int)
	totalWords          int
	categoriesDocuments map[string]int
	categoriesWords     map[string]int
	threshold           float64
}

func NewClassifier(categories []string, threshold float64) Classifier {
	c := Classifier{
		waktu:               make(map[string]map[string]int),
		totalWords:          0,
		categoriesWords:     make(map[string]int),
		categoriesDocuments: make(map[string]int),
		threshold:           threshold,
	}
	for _, category := range categories {
		c.categoriesWords[category] = 0
		c.categoriesDocuments[category] = 0
		c.waktu[category] = make(map[string]int)
	}
	return c
}
func (c *Classifier) addDataTrain(cond Condition) {
	c.datatrain = append(c.datatrain, cond)
}
func (c *Classifier) addDataTest(cond Condition) {
	c.datatest = append(c.datatest, cond)
}

// Train the classifier
func (c *Classifier) Train(category string) {
	c.categoriesWords[category]++
	c.totalWords++
	c.categoriesDocuments[category]++
}

// Classify a document
func (c *Classifier) Classify(dmin float64, dmax float64, tmin float64, tmax float64) (category string) {
	// get all the probabilities of each category
	prob := c.Probabilities(dmin, dmax, tmin, tmax)

	// sort the categories according to probabilities
	var sp []sorted //category, prob
	for c, p := range prob {
		sp = append(sp, sorted{c, p})

		// fmt.Println(sp[1])
	}
	sort.Slice(sp, func(i, j int) bool {
		return sp[i].probability > sp[j].probability
	})

	// if the highest probability is above threshold select that
	if sp[0].probability/sp[1].probability > c.threshold {
		category = sp[0].category
	} else {
		category = "unknown"
	}
	return
}

// Probabilities of each category
func (c *Classifier) Probabilities(dmin float64, dmax float64, tmin float64, tmax float64) (p map[string]float64) {
	p = make(map[string]float64)
	for category := range c.waktu {
		p[category] = c.pCategoryDocument(category, dmin, dmax, tmin, tmax)
	}
	return
}

// p (category | condition1|cond2|cond3|cond4)
func (c *Classifier) pCategoryDocument(category string, dmin float64, dmax float64, tmin float64, tmax float64) float64 {
	class := I_SETOSA
	switch category {
	case "Iris-setosa":
		class = I_SETOSA
		break
	case "Iris-versicolor":
		class = I_VERSICOLOR
		break
	case "Iris-virginica":
		class = I_VIRGINICA
		break
	}
	// fmt.Println(c.pCategory(category))
	return c.pNumericalCategory(class, "l1", dmin) * c.pNumericalCategory(class, "l2", dmax) *
		c.pNumericalCategory(class, "l3", tmin) * c.pNumericalCategory(class, "l4", tmax) *
		c.pCategory(category)
}

// p (condition | category)
func (c *Classifier) pDocumentCategory(category string, condition string) float64 {
	return float64(c.waktu[category][condition]+1) / float64(c.categoriesWords[category])
}

// p (condition numerical | category)
func (c *Classifier) pNumericalCategory(category float64, condition string, cond float64) float64 {
	d := NormalDist{c.avg(condition, category), c.stdev(condition, category)}
	return (d.PDF(cond))
}

// p (category)
func (c *Classifier) pCategory(category string) float64 {
	return float64(c.categoriesDocuments[category]) / float64(len(train))
}
