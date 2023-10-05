package object

func Pipe[T DataTypes](o Object[T], transformers ...Transformer[T]) (Object[T], error) {
	var err error
	for _, t := range transformers {
		o, err = t.Process(o)
		if err != nil {
			return nil, err
		}
	}
	return o, err
}
