{{ define "index/index.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}


<div class="container">
    <h2 class="text-center">一个强大的分布式云盘系统</h2>
    <p class="text-center">本系统是基于Golang,gRPC的分布式云盘系统。</p>
    <hr class="half-rule">
    <div class="row">
        <div class="col-sm-4">
            <div class="row" style="width: 360px; height: 300px;">
                <img src="/static/images/cloud-storage-01.jpg" style="max-height: 300px; max-width: 360px;">
            </div>
            <h3 class="text-center"><b>多端使用</b></h3>
            <p>基于web提供服务，可适配多种消费终端，包括但不限于Windows，Mac，Android，iPhone，iPad等多种平台。</p>
        </div>
        <div class="col-sm-4">
            <div class="row" style="width: 360px; height: 300px;">
                <img src="/static/images/cloud-storage-backup-02.jpg" style="max-height: 300px; max-width: 360px;">
            </div>
            <h3 class="text-center"><b>三重备份</b></h3>
            <p>每个文件将被备份到至少三个存储节点上，集群中有任何一个节点发生宕机都不会导致文件丢失。</p>
        </div>
        <div class="col-sm-4">
            <div class="row" style="width: 360px; height: 300px; vertical-align: middle">
                <img src="/static/images/horizontal-scaling-03.jpg" style="max-height: 300px; max-width: 360px;">
            </div>
            <h3 class="text-center"><b>无限扩展</b></h3>
            <p>空间不够，只需要水平添加集群节点就可以实现空间扩展。理论上可以无限添加节点来达到无限空间的目的。</p>
        </div>
    </div>
    <hr>
    <h2 class="text-center">注册可享受以下优惠权益</h2>
    <div class="row">
        <div class="col-sm-3">
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
                
        </div>
        <div class="col-sm-3">
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
        </div>
        <div class="col-sm-3">
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
        </div>
        <div class="col-sm-3">
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
        </div>
    </div>
    <div class="row text-center">
        <a href="/signup" type="button" class="btn btn-lg btn-primary" style="width: 200px;">立即注册</a>
        &nbsp;&nbsp;&nbsp;&nbsp;
        <a href="/login" type="button" class="btn btn-lg btn-success" style="width: 200px;">登录</a>
    </div>
    <hr>
    <div class="row">
        <p class="text-center">本项目是开源项目，使用Apache协议，开发者可以基于本项目构建自己的应用，也可以向我们贡献自己的代码。<a href="https://github.com/uxff/flexdrive">view source</a></p>
        
    </div>
</div>

{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}
