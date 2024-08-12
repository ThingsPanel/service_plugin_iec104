package form

const ServiceVoucher = `{
	"code": 200,
	"message": "success",
	"data": [
		{
			"dataKey": "address",
			"label": "从站地址",
			"placeholder": "请输入从站地址",
			"type": "input",
			"validate": {
				"message": "从站地址不能为空",
				"required": true,
				"type": "string"
			}
		},
		{
			"dataKey": "port",
			"label": "从站端口",
			"placeholder": "请输入从站端口",
			"type": "input",
			"validate": {
				"message": "从站端口不能为空",
				"required": true,
				"type": "number"
			}
		}
	]
}`
