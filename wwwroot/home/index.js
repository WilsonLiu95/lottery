$(function () {
    function wsInit() {
        var ws = new WebSocket(`ws://${window.location.host}/ws`);
    
        ws.onopen = function () {
            console.log('ws onopen');
            // ws.send('from client: hello');
        };
        ws.onmessage = function (e) {
            console.log('ws onmessage');
            console.log('from server: ' + e.data);
            Lottery.factory(JSON.parse(e.data));
        };
    
        //断线重连
        ws.onclose = function (event) {
            console.log("WebSocket is closed now.");
            wsInit();
        };
    }
    wsInit();

    //backmoney如果值为0的话，那么update 会出现错误，导致动画出错

    var Lottery = {
        data:[],
        userData:[],
        luckyData:[],
        awardArr:[],
        backmoney: 0,
        factory: function(data){
            var type = data['Type'];
            this.data = data;

            switch(type){
                case 'init' :this.init();break;
                case 'ready':this.ready();break;
                case 'start':this.rollScreen();break;
                case 'stop' :this.stop();break;
                case 'pool' :this.addPoolMoney();break;
                case 'show_bag' : this.showBag();break;
                case 'hide_bag' : this.hideBag();break;
                case 'view_on' : this.viewOn();break;
                case 'view_off' : this.viewOff();break;
                case 'view_empty' : this.viewEmpty();break;
            }
        },
        updateBag: function(val){
            val = val || 0;
            updateBag(val);
        },
        updateAwardArr: function(awards) {
            this.awardArr = awards;
        },
        getUserData: function() {
            let self = this;
            $.get('/api/getusers', {}, function (ret) {
                let userDataFormated = ret['data'].map(elem => {
                    return { rtx: elem.Name, name: self.extractChineseName(elem.FullName)};
                });
                self.userData = userDataFormated;
                window.localStorage.setItem('LOTTERY_USER_DATA', JSON.stringify(userDataFormated));
            });
        },
        getAwards: function() {
            let self = this;
            $.ajax({
                url: "/api/getAwards",
                dataType: "json",
                type: "get",
                success: function (info) {
                    if (info.code == 0 && Array.isArray(info.data)) {
                        self.updateAwardArr(info.data);
                    } else {
                        ws.send('from client: getAwards failed.');
                        return 0;
                    }
                },
                error: function (e) {
                    console.log(e);
                }
            });
        },
        //初始化
        init: function(){
            //初始化
            var data = this.data;
            this.backmoney = 0;
            this.updateBag(data['PoolMoney'] || 0);
            toggleMainView(true);
            
            window.localStorage.clear();
            this.getUserData();
            this.getAwards();
        },
        /**
         * 现金池加奖
         */
        addPoolMoney: function(){
            var data = this.data;
            this.updateBag(data['PoolMoney'] || 0);
        },
        viewOn: function(){
            toggleMainView(true);
        },
        viewOff: function(){
            toggleMainView(false);
        },
        viewEmpty: function(){
            hide();
        },
        showBag: function(){
            showBag();
        },
        hideBag: function(){
            hideBag();
        },
        /**
         * ready
         */
        ready: function(){
            console.log('ready');

            this.viewOff();
            
            //初始化奖品信息,根据奖品数量生成人物列表

            var data = this.data,
                prizeNum = data['AllPeopleCount'],
                childRow = this.getRow(prizeNum),
                desc = '';
            this.backmoney = Number(data['BackMoney']) || 0
            this.luckyData = [];//重置
            //重置userData
            this.userData = JSON.parse(window.localStorage.getItem('LOTTERY_USER_DATA'));

            //现金红包判定
            if(Number(data['AwardID']) == 888){
                hideBag();
                if(data['Drawer'] == 'pool'){
                    data['Drawer'] = '';
                }
                ready({
                    packet:{
                        image: '/wwwroot/awards/default.png',
                        typeImg: `/wwwroot/awardtype/0.png`,
                        rtx: data['Drawer'] || '',
                        amount: data['DrawMoney'],//单个金额
                        count: data['AllPeopleCount']//红包个数
                    },
                    childRow:childRow
                });
            }else{
                let award = this.awardArr.find((elem) => {
                    return String(elem.ID) === String(data['AwardID']);
                });
                let typeImg = '0.png';
                if (award && award.hasOwnProperty('ID')) {
                    let typeNum = Number(award['ID']);
                    if (typeNum < 5) {
                        typeImg = '1.png';
                    } else if (typeNum < 9) {
                        typeImg = '2.png';
                    } else if (typeNum < 14) {
                        typeImg = '3.png';
                    } else if (typeNum < 18) {
                        typeImg = '4.png';
                    } else {
                        typeImg = '5.png';
                    }
                }
                let drawer = data['Drawer'] || '';
                ready({
                    gift:{
                        image: `/wwwroot/awards/${award['PicName'] || 'default.png'}`,
                        typeImg: `/wwwroot/awardtype/${typeImg}`,
                        name: drawer + (award['Name'] || '奖品') + (data['Memo'] || ''),
                        desc: award['Description'] || ''
                    },
                    childRow:childRow
                });
            }
        },
        /**
         * 抽奖，屏幕开始滚动
         */
        rollScreen: function(){
            var data = this.data;

            this.luckyData = data['Data'];
            //屏幕滚动
            console.log('roll');
            //调用滚动方法,传入 Lottery.userData即可
            start(Lottery.userData);
        },
        stop: function(){
            console.log('stop');
            luckyData = this.luckyData;

            var data = this.data,
                luckyUser = this.luckyData,
                that    = this;

            //移除已中奖用户
            this.removeLuckyData(luckyUser);
            //动画结束，显示中奖人员
            showResult(this.convertData(luckyUser)).then(function(){
                setTimeout(function(){
                    beforeHide().then(function(){
                        that.updateBag(data['PoolMoney'] || 0);
                    });
                },500);
            });
        },
        /**
         * 移除已中中奖用户,避免重复显示
         */
        removeLuckyData: function(luckyData){

            var i,
                str = JSON.stringify(this.userData);

            //字符串匹配
            for(i = 0 ; i < luckyData.length ; i ++){
                str = str.replace('{"rtx":"' + luckyData[i]['LuckyUserName'] + '","name":"' + this.extractChineseName(luckyData[i]['LuckyUserFullName']) +'"},','');
            }

            window.localStorage.setItem('LOTTERY_USER_DATA',str);
        },
        /**
         * 获取列数
         */
        getRow: function(count){
            var data;
            switch(count){
                case 25: data = [8,9,8];break;
                case 20: data = [6,7,7];break;
                case 10: data = [5,5];break;
                case 9: data = [5,4];break;
                case 8: data = [4,4];break;
                case 4: data = [4];break;
                case 3: data = [3];break;
                case 2: data = [2];break;
                case 1: data = [1];break;
                default : data = this.computeRow(count);
            };

            return data;
        },
        /**
         * 生成奖品列表数组
         */
        computeRow: function(count){
            var col = 0,val = 0,etra = 0,data = [],i;

            if(count < 9){
                col = 1;
                val = count;
            }else{
                if(count == 16){
                    val = 8;
                    col = 2;
                    etra = 0;
                }else{
                    val = 8;
                    col = Math.floor(count/8);
                    etra = count %8;
                }
            }

            for(i = 0 ; i < col; i ++){
                data.push(val);
            }

            etra && data.push(etra);

            return data;
        },
        /**
         * 数据字段转换
         * @param  {[type]} data [description]
         * @return {[type]}      [description]
         */
        convertData: function(data){
            var tempdata = [];
            for(var i = 0 ; i < data.length ; i++){
                tempdata[i] = {};
                tempdata[i]['rtx'] = data[i]['LuckyUserName'].toLowerCase();
                tempdata[i]['avatar'] = this.getHeaderImg(data[i]['LuckyUserName']);
                tempdata[i]['level'] = data[i]['LuckyUserLevel'];
                tempdata[i]['name'] = this.extractChineseName(data[i]['LuckyUserFullName']);
                if (data[i]['LuckyUserLevel'] != 1){
                    tempdata[i]['isLeader'] = true;
                }
                tempdata[i]['backmoney'] = this.backmoney;
            }
            return tempdata;
        },

        /**
         * 从全名获取中文名
         */
        extractChineseName: function (fullname) {
            let rst = fullname.match(/\(([^)]*)\)/);            
            return (Array.isArray(rst) && rst.length > 1) ? rst[1] : '';
        },

        /**
         * 头像规则
         */
        getHeaderImg: function(rtxname){
            return '/wwwroot/images/' + rtxname + '.jpg';
        }
    };

    //首次获取所有用户数据
    (function(){

        var userData = window.localStorage.getItem('LOTTERY_USER_DATA');
        if(userData){
            Lottery.userData = JSON.parse(userData);
        }else{
            Lottery.getUserData();
        }
        window.onload = function(){
            $.get('/api/getpoolmoney',{},function(ret){
                if(ret['code'] == 0){
                    Lottery.updateBag(ret['poolmoney']);
                };
            });

            Lottery.getAwards();
        }

    }());

});
