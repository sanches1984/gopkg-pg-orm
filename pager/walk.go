package pager

// Walk iterates by pages, invoking callback function for each page
func Walk(pageSize int32, fn func(page int32) (int64, error)) error {
	var err error
	var total int64

	p := NewPagerWithPageSize(1, pageSize)
	for {
		total, err = fn(p.GetPage())
		if err != nil {
			return err
		}

		p.SetTotalItems(int32(total))
		err = p.NextPage()
		if LastPageTag.IsTagged(err) {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
