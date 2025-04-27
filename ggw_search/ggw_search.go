package ggw_search

import (
    "fmt"
    //"io/ioutil"
    "os"
    "log"
    "path/filepath"
    "strings"
    "sync"
    "regexp"
    "bufio"
    "archive/zip"
)

var(
    allFiles = make([]string , 0 , 1) ;
    foundPathList map[string]interface{} = make(map[string]interface{}) ;
    
    matchedClauses map[string]interface{}= make(map[string]interface{}); 
    replacedClauses map[string]struct{}= make(map[string]struct{}); 

    regexSpecialChars = `\.+*?()|[]{}^$`;

    globalLock sync.Mutex ;
    result_file *os.File ;
)


func  (s *Search) set_vars(){
    // The code wasn't written structurally at first, and I wasn't in the mood to rewrite it.
    TargetWords = s.TargetWords;
    ReplacementWords = s.ReplacementWords;
    FullSearch = s.FullSearch;
    OnlyGetMatchedClauses = s.OnlyGetMatchedClauses;
    GetMatchedPerFile = s.GetMatchedPerFile;
    IsReplaceAll = s.IsReplaceAll;
    HawMany = s.HawMany;
    StartFromNumber = s.StartFromNumber;
    HasCondition = s.HasCondition;
    MixWithOriginal = s.MixWithOriginal;
    IsModifyMatch = s.IsModifyMatch;
    OutPutMode = s.OutPutMode;
    IsPrintToFile = s.IsPrintToFile;
    ResultFilePath = s.ResultFilePath;
    ValidExtensions = s.ValidExtensions;
    IsStaticDIR = s.IsStaticDIR;
    DIRs = s.DIRs;
    DirToExclude = s.DirToExclude;
    StaticDIRs = s.StaticDIRs;
    RegxToExistBefore = s.RegxToExistBefore;
    RegxToExistAfter = s.RegxToExistAfter;
    InsideCond = s.InsideCond;
    RegxToExistInside = s.RegxToExistInside;
    Regx_NOT_ExistInside = s.Regx_NOT_ExistInside;
    Regx_NOT_ExistBefore = s.Regx_NOT_ExistBefore;
    Regx_NOT_ExistAfter = s.Regx_NOT_ExistAfter;
    ModifyMatch = s.ModifyMatch;
    IsNotRegex = s.IsNotRegex;

    if len(TargetWords) == 0 {
        TargetWords = []string{
            "_all_" , 
        } ; 
    }
    if MixWithOriginal == "" {
        MixWithOriginal = "continue" ;
    }
    if OutPutMode == "" {
        OutPutMode = "println" ;
    }
    if ResultFilePath == "" {
        ResultFilePath = "./ggw_search_result.txt" ;
    }
    if len(DIRs) == 0 {
        DIRs = []string{"./"}; 
    }
}

func (s *Search) Start() {
    s.set_vars();
    var wg sync.WaitGroup
    var closed bool
    var lock sync.Mutex
    
    if !IsStaticDIR {
        for _ , path := range DIRs {
            if _, err := os.Stat(path); os.IsNotExist(err) {
                fmt.Printf("Path does not exist: %s\n", path)
                continue ;
            }

            filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
                path = filepath.ToSlash(path) ;
                for _ , exPath := range DirToExclude {
                    if strings.Contains(path , exPath){
                        return nil ;
                    }
                }

                if err != nil {
                    fmt.Printf("Error accessing path %s: %v\n", path, err)
                    return nil ;
                }
                if !info.IsDir() && hasValidExtension(path, ValidExtensions) {
                            allFiles = append(allFiles, path)
                }          
                return nil ;
            })
        }
    }else{
        allFiles = StaticDIRs ;
    }
    results := make(chan string)
    cunLimit := make(chan struct{} , 300)

    wg.Add(1) 
    go func (results chan string , wg *sync.WaitGroup , closed *bool, lock *sync.Mutex , cunLimit chan struct{} ){  
        for _ , file := range allFiles {
            wg.Add(1)
            cunLimit <- struct{}{} ;
            go search(file, results, wg, closed, lock , cunLimit )
        } 
        wg.Done() ;    
    }(results, &wg, &closed, &lock , cunLimit )

    go func() {
        wg.Wait()
        close(results) 
        close(cunLimit) 
        
    }()

    found := false
    for result := range results {
        if result != "" {
            foundPathList[result] = true;
            found = true
        }
    }
    if !found {
        fmt.Println("Not found")
    }else {
        if FullSearch {
            wg2 := sync.WaitGroup{}
            for path := range foundPathList {
                wg2.Add(1)
                go replace(path,&wg2);
            }
            wg2.Wait() ;
            rCounter := 1 ;
            if len(replacedClauses) > 0 { 
                if IsPrintToFile {
                    var err error ;
                    result_file , err = os.OpenFile(ResultFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666) ;
                    if err != nil {
                        fmt.Println("Error opening file: " , err) ;
                        return ;
                    }
                    defer result_file.Close() ;
                }
                for path := range replacedClauses {
                    if OutPutMode == "println" {
                        if IsPrintToFile {
                            fmt.Fprintln(result_file,rCounter , "- Found in: \"" + filepath.ToSlash(path)  + "\" ; and replace done !")
                        }else{
                            fmt.Println(rCounter , "- Found in: \"" + filepath.ToSlash(path)  + "\" ; and replace done !") ;
                        }
                        rCounter++; 
                    }else if OutPutMode == "showPaths" {
                        fmt.Print(`"`+filepath.ToSlash(path)+`",`) ;
                    }
                }
            }else if !OnlyGetMatchedClauses {
                fmt.Println("nothing found to replace") ;
            }
            if OnlyGetMatchedClauses && OutPutMode != "noOut" {
                if len(matchedClauses) == 0 {
                    fmt.Println("nothing matched") ;
                }
                if !GetMatchedPerFile{
                    if IsPrintToFile {
                        var err error ;
                        result_file , err = os.OpenFile(ResultFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666) ;
                        if err != nil {
                            fmt.Println("Error opening file: " , err) ;
                            return ;
                        }
                        defer result_file.Close() ;
                    }
                    for matched := range matchedClauses {
                        if IsPrintToFile {
                            fmt.Fprintln(result_file,matched)
                        }else{
                            fmt.Println(matched);
                        }
                    }
                }else {
                    i:= 1 ;
                    if IsPrintToFile {
                        var err error ;
                        result_file , err = os.OpenFile(ResultFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666) ;
                        if err != nil {
                            fmt.Println("Error opening file: " , err) ;
                            return ;
                        }
                        defer result_file.Close() ;
                    }
                    for file , matched := range matchedClauses {
                        if IsPrintToFile {
                            fmt.Fprintf(result_file, "path%d - %s :\n", i ,file)
                        }else{
                            fmt.Printf("path%d - %s :\n", i ,file);
                        }
                        for k , m := range matched.([]string) {
                            if IsPrintToFile {
                                fmt.Fprintf(result_file, "    %d - %s\n" , k+1 , strings.TrimSpace(m))
                            }else{
                                fmt.Printf("    %d - %s\n" , k+1 , strings.TrimSpace(m));
                            }
                        }
                        if IsPrintToFile {
                            fmt.Fprintln(result_file)
                        }else{
                            fmt.Println();
                        }
                        i++;
                    }
                }
            }
        }else {
            key := 1 ;
            if IsPrintToFile {
                var err error ;
                result_file , err = os.OpenFile(ResultFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666) ;
                if err != nil {
                    fmt.Println("Error opening file: " , err) ;
                    return ;
                }
                defer result_file.Close() ;
            }
            for path := range foundPathList {
                if IsPrintToFile {
                    fmt.Fprintln(result_file, key , "- Found in: " , path)
                }else{
                    fmt.Println(key , "- Found in: " , path)    ;  
                }
                key ++ ; 
            }
            
        }        
    }
}

func search(filePath string, ch chan string, wg *sync.WaitGroup, closed *bool, lock *sync.Mutex , cunLimit chan struct{}) {
    defer func(){
        <-cunLimit 
        wg.Done()
    }()

    zip,zip_err := zip.OpenReader(filePath);
    fmt.Println();
    fmt.Println("zip err: ", zip_err);
    fmt.Println("zip: ", zip);
    fmt.Println();

    content, err := os.ReadFile(filePath);
    if err != nil {
        fmt.Printf("Error reading file %s: %v\n", filePath, err)
        return
    }

    for _ , TargetWord := range TargetWords {
        TargetWord = strings.TrimSpace(TargetWord);
        if TargetWords[0]== "_all_" || strings.Contains(string(content), TargetWord) {
            lock.Lock()
            if !*closed {
                ch <- filePath
            }
            lock.Unlock()
        }
    }
    content = nil ;
} 

func replace(inputFile string ,wg *sync.WaitGroup) {
    defer wg.Done()
    file, err := os.OpenFile(inputFile, os.O_RDWR, 0644)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer file.Close()
    
	scanner := bufio.NewScanner(file)
	var content string
 
	for scanner.Scan() {
		line := scanner.Text()
		content += line + "\n"
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", inputFile , " , error: " , err)
		return
	}

    newContent := "" ;
    regexMatch(&content , &newContent , inputFile) ;

	file.Seek(0, 0)
	file.Truncate(0)
	writer := bufio.NewWriter(file)
	_, err = writer.Write([]byte(newContent));
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	writer.Flush(); 
}

func regexMatch(content *string, newContent *string , path string){
    for RegxTargetWord , ReplacementWord := range ReplacementWords {
        if IsNotRegex == true {
            RegxTargetWord = escapeRegex(RegxTargetWord);
        }
        re := regexp.MustCompile(RegxTargetWord) ;

        insideRegx := make([]*regexp.Regexp , 0 , 0) ;
        notInsideRegx := make([]*regexp.Regexp , 0 , 0) ;

        if InsideCond {
            for _,regx := range RegxToExistInside {
                if regx == "" {continue};
                insideRegx=append(insideRegx,regexp.MustCompile(regx));
            }
            for _,regx := range Regx_NOT_ExistInside {
                if regx == "" {continue};
                notInsideRegx=append(notInsideRegx,regexp.MustCompile(regx));
            }
        }

        beforeRegx := regexp.MustCompile(RegxToExistBefore);
        afterRegx := regexp.MustCompile(RegxToExistAfter);
        notBeforeRegx := regexp.MustCompile(Regx_NOT_ExistBefore);
        notAfterRegx := regexp.MustCompile(Regx_NOT_ExistAfter); 
        
        matchCount := 0 ;
        repCount := 0 ;

        if IsReplaceAll { 
            *newContent = re.ReplaceAllStringFunc(*content, func(match string) string {
                modified := "" ;
                if IsModifyMatch {
                    modified = ModifyMatch(match , ReplacementWord);
                }
                if IsModifyMatch && modified == "" {
                    return match ;
                }else if modified != ""{
                    match = modified ;
                }

                if OnlyGetMatchedClauses == true {
                    addMatchedClauses(match , path);
                    return match ;
                }
                return repMatched(match, path , ReplacementWord) ;
            }) ;
        }else{
            matchedIndexes := re.FindAllStringIndex(*content , -1);
            currentIndex := make([]int ,0, 1) ;

            contentBetween := "" ;
            if InsideCond {
                indexOfBeforeRegx := beforeRegx.FindAllStringIndex(*content , -1) ;
                indexOfAfterRegx := afterRegx.FindAllStringIndex(*content , -1) ;
                if len(indexOfBeforeRegx) > 2 || len(indexOfAfterRegx) > 2 {
                    fmt.Println("Warning: indexOfBeforeRegx || indexOfAfterRegx > 2");
                }
                if indexOfBeforeRegx != nil && indexOfAfterRegx != nil &&
                len(indexOfBeforeRegx) > 0 && len(indexOfAfterRegx) > 0{
                    if RegxToExistBefore == RegxToExistAfter {
                        contentBetween = (*content)[indexOfBeforeRegx[0][1]:indexOfAfterRegx[1][0]] ;
                    }else{
                        contentBetween = (*content)[indexOfBeforeRegx[0][1]:indexOfAfterRegx[0][0]] ;
                    }
                }
            }

            *newContent = re.ReplaceAllStringFunc(*content, func(match string) string {
                matchCount ++ ;
                if HawMany == -1 || matchCount >= StartFromNumber {
                    if HawMany == -1 || repCount < HawMany {
                        if !HasCondition {

                            modified := "" ;
                            if IsModifyMatch {
                                modified = ModifyMatch(match , ReplacementWord);
                            }
                            if IsModifyMatch && modified == "" {
                                return match ;
                            }else if modified != ""{
                                match = modified ;
                            }

                            if OnlyGetMatchedClauses == true {
                                addMatchedClauses(match , path);
                                return match ;
                            }
                            repCount++;
                            return repMatched(match, path , ReplacementWord) ;
                        }

                        
                        currentIndex = matchedIndexes[matchCount-1] ;

                        beforePart := (*content)[0:currentIndex[0]] ;
                        afterPart := (*content)[currentIndex[1]:] ;

                        existsBefore := RegxToExistBefore == "" || beforeRegx.MatchString(beforePart) ;
                        existsBetween := true ; 
                        if contentBetween != "" && InsideCond {
                            for _ , regx := range insideRegx {
                                if !regx.MatchString(contentBetween) {
                                    existsBetween = false ;
                                    break;
                                }
                            }
                            for _ , regx := range notInsideRegx {
                                if regx.MatchString(contentBetween) {
                                    existsBetween = false ;
                                    break;
                                }
                            }
                        }

                        existsAfter := RegxToExistAfter == "" || afterRegx.MatchString(afterPart) ;
                        
                        notExistBefore := Regx_NOT_ExistBefore == "" || !notBeforeRegx.MatchString(beforePart) ;
                        notExistAfter := Regx_NOT_ExistAfter == "" || !notAfterRegx.MatchString(afterPart) ;
                        
                        if existsBefore && existsAfter && notExistBefore && notExistAfter && existsBetween {

                            modified := "" ;
                            if IsModifyMatch {
                                modified = ModifyMatch(match , ReplacementWord);
                            }
                            if IsModifyMatch && modified == "" {
                                return match ;
                            }else if modified != ""{
                                match = modified ;
                            }
                            
                            if OnlyGetMatchedClauses == true {
                                addMatchedClauses(match , path);
                                return match ;
                            }
                            repCount++ ;
                            return repMatched(match, path , ReplacementWord) ;
                        } else {
                            return match ;
                        }

                    }else{
                        return match ;
                    }
                }else{
                    return match ;
                }
                log.Fatal("error in ReplaceAllStringFunc \n");
                return "" ;
            }) ;
        }
        *content = *newContent ;
    }
}

func repMatched(match , path , ReplacementWord string) string{
    globalLock.Lock();
    replacedClauses[path] = struct{}{} ;
    globalLock.Unlock();
    switch MixWithOriginal {
        case "replace": return ReplacementWord ;
        case "after": return match + ReplacementWord ;
        case "before": return ReplacementWord + match ;
        case "continue": return match ;
    }
    log.Fatal("error in repMatched \n");
    return "" ;
}

func addMatchedClauses(match , path string){
    globalLock.Lock();
    if !GetMatchedPerFile { 
        matchedClauses[match]=nil;
    }else{
        tempArr , tempArrExist := matchedClauses[path].([]string) ;
        if !tempArrExist {
            tempArr = make([]string , 0 , 1) ;
            matchedClauses[path] = tempArr ;
        }
        tempArr = append(tempArr , match) ;
        matchedClauses[path] = tempArr ;
    }
    globalLock.Unlock(); 
}

func hasValidExtension(path string, ValidExtensions []string) bool {
	for _, ext := range ValidExtensions {
        if ext == "all" {return true}
		if strings.HasSuffix(strings.ToLower(path), strings.ToLower(ext)) {
			return true
		}
	}
	return false
}

func escapeRegex(input string) string {
	var escaped strings.Builder
	for _, ch := range input {
		if strings.ContainsRune(regexSpecialChars, ch) {
			escaped.WriteRune('\\')
		}
		escaped.WriteRune(ch)
	}
	return escaped.String()
}