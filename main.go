package main

import (
	"github.com/gogetwebsite/ggw_search";
	"fmt";
)

func main() {
	s := ggw_search.Search{};

	s.TargetWords = []string{
		//"_all_",
		"Hello from txt4",
	};
	s.ReplacementWords = map[string]string{
		`Hello from txt4`: "",
	};
	s.IsStaticDIR  = true;
	s.StaticDIRs = []string{"./test/a2.docx" , "./test/a.txt"};
	s.FullSearch = false; 
	s.IsReplaceAll = true;
	s.OnlyGetMatchedClauses = true; // only simulate the search and dont't replace anything
	s.GetMatchedPerFile = true;
	s.ValidExtensions = []string{".docx"};
	s.DIRs = []string{"./"};
	s.IsModifyMatch = true;
	s.ModifyMatch = func(match , replacementWord string)string{
		fmt.Println("matched: " , match);
		return match;
	}
	s.Start();
}
