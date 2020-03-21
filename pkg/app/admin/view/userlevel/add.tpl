{{ define "userlevel/add.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}

<div class="container">
    <div class="row">
        <ul class="breadcrumb">
            <li><a href="/">首页</a></li>
            <li><a href="/userlevel/list">会员升级包</a></li>
            <li class="active">添加会员升级包</li>
        </ul>
        <div class="col-md-8 col-md-offset-2">
    		<div class="panel panel-default">
			  	<div class="panel-heading text-center">
			    	<h3 class="panel-title"><strong>
                        添加会员升级包
                    </strong></h3>
			 	</div> 

			  	<div class="panel-body">
                    <form accept-charset="utf-8" role="form" class="form-horizontal" method="POST" action="/userlevel/add">
                        <div class="form-group">
                            <label class="col-sm-3 control-label">升级包名称</label>
                            <div class="col-sm-8">
                                <input class="form-control" placeholder="例: 黄金会员" name="name" value="" type="text" required />
                            </div>
                        </div>
                        <div class="form-group">
                            <label class="col-sm-3 control-label">介绍</label>
                            <div class="col-sm-8">
                                <input class="form-control" placeholder="例: 送离线下载" name="desc" value="" type="text" />
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="" class="col-sm-3 control-label">配额空间(KB)</label>
                            <div class="col-sm-8">
                                <input class="form-control" name="quotaSpace" value="0" type="number" required/>
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="" class="col-sm-3 control-label">付款价格(分)</label>
                            <div class="col-sm-8">
                                <input class="form-control" placeholder="0.01" name="price" value="0" type="number" required />
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="" class="col-sm-3 control-label">展示原价(分)</label>
                            <div class="col-sm-8">
                                <input class="form-control" placeholder="0.01" name="primeCost" value="0" type="number" />
                            </div>
                        </div>
                        <div class="form-group">
                            <label class="col-sm-3 control-label">是否默认升级包</label>
                            <div class="col-sm-8">
                                <input class="form-control" name="isDefault" value="1" type="checkbox" />
                            </div>
                        </div>
                                
                        <div class="form-group text-center">
                            <div class="col-sm-5"></div>
                            <div class="col-sm-2">
                                <input class="btn btn-success btn-block" type="submit" value="提交">
                            </div>
                            <div class="col-sm-5"></div>
                            </div>
                        </form>
                </div>
			</div>
		</div>
	</div>
</div>
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}