
build:
	go build -o step0_repl Step0_repl/step0_repl.go
	go build -o step1_read_print Step1_read_print/step1_read_print.go
	go build -o step2_eval Step2_eval/step2_eval.go
	mv step0_repl ../
	mv step1_read_print ../
	mv step2_eval ../
Run:
	go build -o step0_repl Step0_repl/step0_repl.go
	go build -o step1_read_print Step1_read_print/step1_read_print.go
	go build -o step2_eval Step2_eval/step2_eval.go
	./step0_repl
clean:
	go clean
	rm step0_repl/step0_repl