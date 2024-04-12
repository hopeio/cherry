package cron

import (
	"github.com/robfig/cron/v3"
)

/*
cron表达式代表一个时间的集合，使用6个空格分隔的字段表示：

字段名	是否必须	允许的值　	允许的特定字符
秒(Seconds)	是	0-59	* / , -
分(Minute)	是	0-59	* / , -
时(Hours)	是	0-23	* / , -
日(Day of month)	是	1-31	* / , - ?
月(Month)	是	1-12 或 JAN-DEC	* / , -
星期(Day of week)	否	0-6 或 SUM-SAT	* / , - ?

　　　　注：

　　　　1.月(Month)和星期(Day of week)字段的值不区分大小写，如：SUN、Sun 和 sun 是一样的。

　　　　2.星期(Day of week)字段如果没提供，相当于是 *

复制代码
# ┌───────────── min (0 - 59)
# │ ┌────────────── hour (0 - 23)
# │ │ ┌─────────────── day of month (1 - 31)
# │ │ │ ┌──────────────── month (1 - 12)
# │ │ │ │ ┌───────────────── day of week (0 - 6) (0 to 6 are Sunday to
# │ │ │ │ │                  Saturday, or use names; 7 is also Sunday)
# │ │ │ │ │
# │ │ │ │ │
# * * * * *  command to execute
　1）星号(*)
　　　　表示 cron 表达式能匹配该字段的所有值。如在第5个字段使用星号(month)，表示每个月

　　2）斜线(/)
　　　　表示增长间隔，如第1个字段(minutes) 值是 3-59/15，表示每小时的第3分钟开始执行一次，之后每隔 15 分钟执行一次（即 3、18、33、48 这些时间点执行），这里也可以表示为：3/15

　　3）逗号(,)
　　　　用于枚举值，如第6个字段值是 MON,WED,FRI，表示 星期一、三、五 执行

　　4）连字号(-)
　　　　表示一个范围，如第3个字段的值为 9-17 表示 9am 到 5pm 直接每个小时（包括9和17）

　　5）问号(?)
　　　　只用于 日(Day of month) 和 星期(Day of week)，表示不指定值，可以用于代替 *

　　6）L，W，#
　　　　Go中没有L，W，#的用法，下文作解释。
*/

//cron举例说明
//　　　　每隔5秒执行一次：*/5 * * * * ?
//
//每隔1分钟执行一次：0 */1 * * * ?
//
//每天23点执行一次：0 0 23 * * ?
//
//每天凌晨1点执行一次：0 0 1 * * ?
//
//每月1号凌晨1点执行一次：0 0 1 1 * ?
//
//在26分、29分、33分执行一次：0 26,29,33 * * * ?
//
//每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?

var funcMap = map[string]func(){}

var jobMap = map[string]cron.Job{}

func init() {
	funcMap["0 10 10 * * *"] = Job
	jobMap["0 0 2 * * *"] = RedisTo{}
}

// New 构造cron
func New() *cron.Cron {
	c := cron.New()
	for spec, cmd := range funcMap {
		c.AddFunc(spec, cmd)
	}
	for spec, cmd := range jobMap {
		c.AddJob(spec, cmd)
	}
	return c
}
