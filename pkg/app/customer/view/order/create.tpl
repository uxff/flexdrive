{{ define "order/create.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}



<div class="container">
<form id="mockpay" method="POST" action="/my/order/create">

    <h2 class="text-center">购买升级享受以下优惠权益</h2>
    <div class="row">
        <div class="col-sm-3 levels">
            <label for="level1" style="display: block;padding: 8px;">
            <div class="panel panel-success">
                <div class="panel-heading">
                    <h3 class="text-center">青铜会员</h3>
                </div>
                <div class="panel-body">
                    <ul style="margin-left: 10%">
                        <li>空间：500MB</li>
                        <li>不支持离线下载</li>
                    </ul>
                    <p class="text-center"><del>原价格：1元</del></p>
                    <p class="text-success text-center"><b>现价格：免费 注册即送</b></p>
                </div>
            </div>
            <input type="radio" name="level" id="level1" value="1"></label>
        </div>
        <div class="col-sm-3 levels">
            <label for="level2" style="display: block;padding: 8px;">
            <div class="panel panel-info">
                <div class="panel-heading">
                    <h3 class="text-center">白银会员</h3>
                </div>
                <div class="panel-body">
                    <ul style="margin-left: 10%">
                        <li>空间：5GB</li>
                        <li class="active ">支持离线下载</li>
                    </ul>
                    <p class="text-center"><del>原价格：5元</del></p>
                    <p class="text-success text-center"><b>现价格：3元</b></p>
                </div>
            </div>
            <input type="radio" name="level" id="level2" value="2"></label>
        </div>
        <div class="col-sm-3 levels">
            <label for="level3" style="display: block;padding: 8px;">
            <div class="panel panel-warning">
                <div class="panel-heading">
                    <h3 class="text-center">黄金会员</h3>
                </div>
                <div class="panel-body">
                    <ul style="margin-left: 10%">
                        <li>空间：50GB</li>
                        <li class="active">支持离线下载</li>
                    </ul>
                    <p class="text-center"><del>原价格：50元</del></p>
                    <p class="text-success text-center"><b>现价格：10元</b></p>
                </div>
            </div>
            <input type="radio" name="level" id="level3" value="3"></label>
        </div>
        <div class="col-sm-3 levels">
            <label for="level4" style="display: block;padding: 8px;">
            <div class="panel panel-primary">
                <div class="panel-heading">
                    <h3 class="text-center">钻石会员</h3>
                </div>
                <div class="panel-body">
                    <ul style="margin-left: 10%">
                        <li>空间：500GB</li>
                        <li class="active">支持离线下载</li>
                    </ul>
                    <p class="text-center"><del>原价格：500元</del></p>
                    <p class="text-success text-center"><b>现价格：50元</b></p>
                </div>
            </div>
            <input type="radio" name="level" id="level4" value="4"></label>
        </div>
    </div>
    <div class="row text-center">
        <input type="submit" class="btn btn-lg btn-primary" style="width: 200px;">
    </div>
    <hr>
    <div class="row">
        <p class="text-center"></p>
        
    </div>
</form>
</div>


<script type="text/javascript">

$(".levels").on('click', function(){
    // $('.levels').removeClass('purple');
    // $(this).addClass('purple');
    $('.levels').css({background: '#ffffff'});
    $(this).css({background: '#9076c3'});
});


</script>
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}
