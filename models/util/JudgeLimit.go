package util

func TotalPage(total int, pagesize int) int{
	// totalPage = 0
	totalPage := 0
	if total % pagesize == 0 {
		totalPage = total / pagesize
	} else {
		totalPage = total / pagesize + 1
	}
	return totalPage
}

func JudgePage(total int,pagenum int,pagesize int) (int,int,int){
	totalPage := 0
	if pagenum == 0{
		if pagesize == 0 {
			pagesize = 10
			pagenum = 1
			totalPage = TotalPage(total, pagesize)
		}else {
			if pagesize > 0 && pagesize < 100{
				totalPage = TotalPage(total, pagesize)
				pagenum = 1
			}else {
				pagesize = 10
				totalPage = TotalPage(total, pagesize)
				pagenum = 1
			}
		}
	}else {
		if pagesize == 0 {
			pagesize = 10
			totalPage = TotalPage(total, pagesize)
			if pagenum < 0 || pagenum > totalPage{
				pagenum = 1
			}
		}else if pagesize > 0 && pagesize < 100{
			totalPage = TotalPage(total, pagesize)
			if pagenum < 0 || pagenum > totalPage {
				pagenum = 1
			}
		}else {
			pagesize = 10
			totalPage = TotalPage(total, pagesize)
			if pagenum < 0 || pagenum > totalPage {
				pagenum = 1
			}
		}
	}
	return totalPage,pagenum,pagesize
}
