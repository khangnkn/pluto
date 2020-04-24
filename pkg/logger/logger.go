package logger

import "go.uber.org/zap"

var lg *zap.SugaredLogger

func Initlialize(prod bool) {
	if prod {
		l, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		lg = l.Sugar()
		lg.Infof("registered production logger")
	} else {
		l, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		lg = l.Sugar()
		lg.Infof("registered development logger")
	}
}

func Info(args ...interface{}) {
	lg.Info(args)
}

func Infof(template string, args ...interface{}) {
	lg.Infof(template, args)
}

func Error(args ...interface{}) {
	lg.Error(args)
}

func Panic(args ...interface{}) {
	lg.Panic(args)
}
