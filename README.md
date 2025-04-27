# 🔍 Fast Find and Replace

A super fast tool to search and (optionally) replace text across tons of files!

## 🚀 What it does

- Opens files **concurrently**  
- **Searches** for specific text/regexp inside them
- **Replaces** that text (only if you want)  
- Uses your CPU power to the max for blazing speed ⚡

## ⚙️ Configurable

You’re in control:

- Set what text to search for  
- What to replace it with (or not)  
- Whether to actually replace or just search  
- Which folder(s) to scan  
- Which specific file(s) to scan  
- Which file extensions to include
- Which paths to exclude
- Whether a conditional regexp exists before, after or inside the specified text.
- Perform a function for each matched text.
- Print the results to stdout or in a file.
- Specify how many matches should be replaced and which one to start from.

## 💻 How to run

main.go :

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
    s.StaticDIRs = []string{"./a.txt"};
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

## 🧠 Why use it?

- Way faster than basic scripts  
- Handles thousands of files at once  
- Fully customizable  
- Open source and easy to hack on
- Easy replace your texts in the codes with conditions

## 📝 Full config
type Search struct{
    TargetWords []string; // Sifting through data, to instantly search and add the path to the searchable path list.
    ReplacementWords map[string]string; // `(Regex)TargetWord` : "to raplace with" ,
    IsNotRegex bool; // TargetWord is not a regex
    FullSearch bool; // use 'ReplacementWords' keys to search in all lines (and replace or not)
    OnlyGetMatchedClauses bool; // to get all matched items without any replacement
    GetMatchedPerFile bool; // to show matched clauses separately in each file, or only return one as an instance
    IsReplaceAll bool; // Should be set to true to complete the search based on 'ReplacementWords' keys. handle no replacement by OnlyGetMatchedClauses
    MixWithOriginal string; // replacement mode: "replace","after","before","continue"(skip, for testing or ...)
    OutPutMode string; // "println","showPaths","noOut"
    IsPrintToFile bool; // print to file (for big outputs), or print to stdout
    ResultFilePath string; // file path to print to file 
    ValidExtensions []string; // set {"all"} to serach in all files
    IsStaticDIR bool; // search a whole directory, or specific paths?
    StaticDIRs []string; // specify your Static DIRs to avoid dynamic recursive search
    DIRs []string; // specify directories to search dynamic recursive search
    DirToExclude []string ; // specify part of directories to Exclude searching
    HasCondition bool; // whether replacement has condition, if yes, IsReplaceAll should be set to false
    HawMany int; // for conditional replace, set it -1 to do for all matches
    StartFromNumber int; // for conditional replace, replace from which matched item?
    RegxToExistBefore string; // for conditional replace
    RegxToExistAfter string; // for conditional replace
    InsideCond bool; // for conditional replace
    RegxToExistInside []string; // for conditional replace
    Regx_NOT_ExistInside []string; // for conditional replace
    Regx_NOT_ExistBefore string; // for conditional replace
    Regx_NOT_ExistAfter string; // for conditional replace
    IsModifyMatch bool; // to run your custom code (the ModifyMatch function) for each matched
    ModifyMatch func(match , replacementWord string) string; // run your custom code for each matched
} 
