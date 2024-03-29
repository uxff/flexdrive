{{ define "userlevel/edit.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}

<div class="container">
    <div class="row">
        <ul class="breadcrumb">
          <li><a href="/">首页</a></li>
          <li><a href="/userlevel/list">会员升级包</a></li>
          <li class="active">编辑会员升级包</li>
        </ul>
    	<div class="col-md-8 col-md-offset-2">
    		<div class="panel panel-default">
			  	<div class="panel-heading text-center">
			    	<h3 class="panel-title"><strong>编辑会员升级包</strong></h3>
			 	</div> 

            <div class="panel-body">
                <form accept-charset="utf-8" role="form" class="form-horizontal" method="POST" action="/userlevel/edit/{{.levelEnt.Id}}">
                    <div class="form-group">
                        <label for="" class="col-sm-3 control-label">会员升级包id</label>
                        <div class="col-sm-8 control-label" style="text-align: left">
                            {{.levelEnt.Id}}
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-3 control-label">升级包名称</label>
                        <div class="col-sm-8">
                            <input class="form-control" placeholder="例: 运营" name="name" value="{{.levelEnt.Name}}" type="input" required />
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-3 control-label">介绍</label>
                        <div class="col-sm-8">
                            <input class="form-control" placeholder="例: 送某某礼品" name="name" value="{{.levelEnt.Desc}}" type="input"  />
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-3 control-label">配额空间(KB)</label>
                        <div class="col-sm-8">
                            <input class="form-control" placeholder="0" name="quotaSpace" value="{{.levelEnt.QuotaSpace}}" type="number" required />
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-3 control-label">购买价格(分)</label>
                        <div class="col-sm-8">
                            <input class="form-control" name="price" value="{{.levelEnt.Price}}" type="number" required />
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-3 control-label">展示原价(分)</label>
                        <div class="col-sm-8">
                            <input class="form-control" name="primeCost" value="{{.levelEnt.PrimeCost}}" type="number" />
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-3 control-label">是否默认升级包</label>
                        <div class="col-sm-8">
                            <input class="form-control" name="isDefault" {{if eq .levelEnt.IsDefault 1}}checked{{end}} value="1" type="checkbox" />
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