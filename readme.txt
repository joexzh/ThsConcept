http://basic.10jqka.com.cn/603737/concept.html
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"
`<h1 style="margin:3px 0px 0px 0px">\s*\d{6}\s*</h1>`
`cid="(\d*)"`




1   股票列表文件
    http://www.shdjt.com/js/lib/astock.js
    regexp=`~(?:[a-z]{2})?([0-9]*)` matches group 1
    用 Regester 测试正则表达式, E:\Program Files\Regester

2   股票列表筛选
    True
        0开头
        6开头
        3开头
    False
        39开头
    筛选完后去重

3   判断页面是否合法
    http://basic.10jqka.com.cn/166007/concept.html
    regexp=`cid="(\d*)"` matches group 1, 即概念id
    去重

4   然后调用同花顺的概念接口, rap2.taobao.org
    每个概念存入ths_stock数据库, 概念(一)--stock(多),

    table concept
        id: "" primary key
        name: ""
        updateDate: int64

    table concept_stock

        conceptId
        stockCode
        concept_description: ""

    table concept_stock
        stockCode varchar(6)
        stockName


    db collection 格式:
    {
        conceptId: ""
        conceptName: ""
        updateDate: int64 // 更新时间
        stocks: [
            {
                stockCode: ""
                stockName: ""
                description: "个股概念描述"
            }
        ]
    }

5   定时每天晚上10:20爬一次数据
    // TODO 完成后微信提醒

6   搭建web服务
    搜 股票名称 或 股票代码, 出来对应的概念列表, group
    {
        stockCode: ""
        stockName: ""
        concepts: [
            {
                conceptId: ""
                conceptName: ""
                reportDate: int64
                description: ""
            }
        ]
    }

    搜 概念名称, 支持模糊查询, 出来对应的股票列表, group
    {
        conceptId: ""
        conceptName: "模糊查询"
        reportDate: int64
        stocks: [
            {
                stockCode: ""
                stockName: ""
                description: ""
            }
        ]
    }


每只股票的概念, 时间
{
    conceptId: "",
    stockCode: "",
    stockName: "",
    description: "",
    lastModified: 16124521342
}