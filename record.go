package ddns

// Record 保存记录信息
type Record struct {
	Type       string
	DomainName string
	RecordID   string
	Value      string
	Prefix     string
}

func (r *Record) Merge(diff *Record) {
	if len(diff.Type) != 0 {
		r.Type = diff.Type
	}
	if len(diff.Value) == 0 {
		r.Value = diff.Value
	}
}

// 查找修改方案,返回数组有三个元素，第一个是create，第二个是delete，第三个是update
func findBestSolution(rcds []*Record, tp, domain, prefix string, addrs []string) (del, crt, upd []*Record, err error) {
	del = make([]*Record, 0, len(rcds))
	crt = make([]*Record, 0, len(rcds))
	upd = make([]*Record, 0, len(rcds))
	// dirty用来标记是否需要修改
	// 第一位用来标记类型是否被修改 第二位用来标记源ip是否有效 第三位用来标记ip是否被修改
	var dirty = make(map[*Record]int8)
	var status int8
	for _, rcd := range rcds {
		status = 0
		// 如果记录类型不同，修改记录类型
		if rcd.Type != tp {
			rcd.Type = tp
			status = status | 1
		}
		// 如果地址没有出现过，那么就记录他为空
		for i, v := range addrs {
			if rcd.Value == v {
				status = status | 2
				// 从ip中删除掉
				addrs = append(addrs[:i], addrs[i+1:]...)
				break
			}
		}
		dirty[rcd] = status
	}
	// 最后，把没有匹配成功的ip放入到无效的记录中
	for rcd, v := range dirty {
		// 如果所有记录都匹配完成了，直接标记无效的ip记录为删除
		if len(addrs) == 0 && v&2 == 0 {
			del = append(del, rcd)
			continue
		}
		// 寻找没有匹配成功的记录,把它的记录值修改
		if v&2 == 0 {
			rcd.Value = addrs[0]
			upd = append(upd, rcd)
			addrs = addrs[1:]
		}
	}
	// 如果地址还有剩下的，就直接创建新记录
	for _, v := range addrs {
		crt = append(crt, &Record{
			DomainName: domain,
			Type:       tp,
			Prefix:     prefix,
			Value:      v,
		})
	}
	return crt, del, upd, nil
}
