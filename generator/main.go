package main

func main() {
	for _, p := range packages {
		Generate(p, cmds)
	}
}
