package form

const DeviceVoucher = `{
    "code": 200,
    "message": "success",
    "data": [
        {
            "dataKey": "address",
            "label": "数据库地址",
            "placeholder": "请输入数据库地址",
            "type": "input",
            "validate": {
                "message": "数据库地址不能为空",
                "required": true,
                "type": "string"
            }
        },
        {
            "dataKey": "password",
            "label": "数据库密码",
            "placeholder": "请输入数据库密码",
            "type": "input",
            "validate": {
                "message": "",
                "required": false,
                "type": "string"
            }
        },
        {
            "dataKey": "interval",
            "label": "轮训间隔(秒)",
            "placeholder": "请输入轮训间隔",
            "type": "input",
            "validate": {
                "message": "",
                "required": false,
                "type": "number"
            }
        }
    ]
}`
