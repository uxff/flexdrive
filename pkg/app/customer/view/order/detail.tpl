{{ define "order/detail.tpl" }}
{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}
<div class="container">

    <div class="row">
        <ul class="breadcrumb">
            <li><a href="/">首页</a></li>
            <li><a href="/my/order/list">我的订单</a></li>
            <li class="active">订单详情</li>
        </ul>
    </div>

    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-8" >
            <div class="panel panel-success">
                <div class="panel-heading">
                    <h3 class="text-center">订单号： {{.Order.Id}}</h3>
                </div>
                <div class="panel-body">
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-4" >
                            下单时间
                        </div>
                        <div class="col-md-4" >
                            {{.Order.Created}}
                        </div>
                        <div class="col-md-2"></div>
                    </div>
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-4" >
                            会员账号
                        </div>
                        <div class="col-md-4" >
                            {{.LoginInfo.UserEnt.Email}}
                        </div>
                        <div class="col-md-2"></div>
                    </div>
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-4" >
                            购买内容
                        </div>
                        <div class="col-md-4" >
                            {{.Level.Name}}({{.Level.Id}})
                        </div>
                        <div class="col-md-2"></div>
                    </div>
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-4" >
                            获得权益
                        </div>
                        <div class="col-md-4" >
                            空间增加 {{space4Human .Order.AwardSpace}}
                        </div>
                        <div class="col-md-2"></div>
                    </div>
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-4" >
                            金额
                        </div>
                        <div class="col-md-4" >
                            {{amount4Human .Order.TotalAmount}} 元
                        </div>
                        <div class="col-md-2"></div>
                    </div>
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-4" >
                            状态
                        </div>
                        <div class="col-md-4" >
                            {{orderStatus .Order.Status}}
                        </div>
                        <div class="col-md-2"></div>
                    </div>

                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-8 text-center" >
                            {{if eq .Order.Status 1}}
                            <a href="/my/order/mockpay/{{.Order.Id}}" class="btn btn-success" target="_blank">去付款</a>
                            {{end}}
                        </div>
                        <div class="col-md-2"></div>
                    </div>

                </div>
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>
</div>
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}
{{end}}