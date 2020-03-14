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
    <div class="row">
        <p class="text-center">本项目是开源项目，使用Apache协议，开发者可以基于本项目构建自己的应用，也可以向我们贡献自己的代码。<a href="https://github.com/uxff/flexdrive">view source</a></p>
        
    </div>
</div>

{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}
