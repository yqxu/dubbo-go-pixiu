package tester

type Tester struct {
	operator string
	_type    string // batch / single
	priority int
	lazy     bool //if lazy load or not
}

type batchTester struct {
	Tester
	valueSet interface{}
}

type singleTester struct {
	Tester
	staticValue interface{}
}

func (t *Tester) IsMatch(tested interface{}) {

}
