package samplepkg

type Input struct {
	Verbose bool
}

type Output struct {
}

func SampleFunc(in *Input) (*Output, error) {
	return &Output{}, nil
}
