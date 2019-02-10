package redos

var fuzzStrings = []string{
	// https://github.com/rust-lang/regex/blob/master/examples/shootout-regex-dna-cheat.rs
	"agggtaaa|tttaccct",
	"[cgt]gggtaaa|tttaccc[acg]",
	"a[act]ggtaaa|tttacc[agt]t",
	"ag[act]gtaaa|tttac[agt]ct",
	"agg[act]taaa|ttta[agt]cct",
	"aggg[acg]aaa|ttt[cgt]ccct",
	"agggt[cgt]aa|tt[acg]accct",
	"agggta[cgt]a|t[acg]taccct",
	"agggtaa[cgt]|[acg]ttaccct",
	// https://www.owasp.org/index.php/Regular_expression_Denial_of_Service_-_ReDoS
	// Examples of Evil Patterns:
	// (a+)+
	// ([a-zA-Z]+)*
	// (a|aa)+
	// (a|a?)+
	// (.*a){x} | for x > 10
	// ^([a-zA-Z0-9])(([\-.]|[_]+)?([a-zA-Z0-9]+))*(@){1}[a-z0-9]+[.]{1}(([a-z]{2,3})|([a-z]{2,3}[.]{1}[a-z]{2,3}))$
	// ^(([a-z])+.)+[A-Z]([a-z])+$
	// All the above are susceptible to the input aaaaaaaaaaaaaaaaaaaaaaaa!
	// (The minimum input length might change slightly, when using faster or slower machines).
	"aaaaaaaaaaaaaaaaaaaaaaaa!",
	"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa!",
}
