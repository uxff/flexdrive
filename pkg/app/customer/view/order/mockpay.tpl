{{ define "order/mockpay.tpl" }}

<form id="mockpay" method="POST" action="/my/order/notify">
    <p>
        本页面用于模拟本系统的订单支付跳转到第三方支付页面
    </p>
    <p>
        请在1分钟内支付，否则订单会失效，失效后请重新打开支付
    </p>
    <table style="border: 1px">
        <tr>
            <td>订单号</td>
            <td>{{.Order.Id}}</td>
        </tr>
        <tr>
            <td>用户Id</td>
            <td><input type="text" name="userId" readonly value="{{.Order.UserId}}"></td>
        </tr>
        <tr>
            <td>用户Email</td>
            <td>{{.User.Email}}</td>
        </tr>
        <tr>
            <td>订单号</td>
            <td><input type="text" name="orderId" readonly value="{{.Order.Id}}"></td>
        </tr>
        <tr>
            <td>订单内容</td>
            <td>{{.Order.LevelName}}</td>
        </tr>
        <tr>
            <td>订单金额</td>
            <td>{{.Order.TotalAmount}}</td>
        </tr>
        <tr>
            <td>支付机构订单号</td>
            <td><input type="text" name="outOrderNo" readonly value="{{.OutOrderNo}}"></td>
        </tr>
        <tr>
            <td>输入支付验证码({{.VerifyCode}})</td>
            <td><input type="text" name="verifyCode"></td>
        </tr>
    </table>
    <input type="hidden" name="token" value="{{.Token}}">
    <button type="submit">确认支付</button>
    <a href="javascript:history.go(-1);">返回</a>
</form>
{{end}}