package sample

import (
	"errors"
	"fmt"
	"qiyuesuo/sdk/http"
	"qiyuesuo/sdk/model"
	"qiyuesuo/sdk/request"
	"testing"
)

/*
签署示例：该场景模拟一个人事合同的场景，即平台方公司与员工签署合同，平台方公司先签署，员工后签。
具体操作：
1、提前在登录契约锁云平台，配置好业务分类，并指定好文件模板、签署流程、签署时使用的公章、签署位置等信息；
2、通过调用接口，指定好签署方和模板参数，即可发起合同；
3、通过接口签署本方企业，再通过接口获取第三方个人的签署链接，个人打开签署链接签署合同。
*/
func TestSendByCategory(t *testing.T) {
	//sdkClient := http.NewSdkClient("https://openapi.qiyuesuo.cn", "tBYw111111", "yHnVOk11MNUT97BKQnN5xnul111111")
	sdkClient := http.NewSdkClient("更换为开放平台请求地址", "更换为您开放平台 App Token", "更换为您开放平台App Secret")
	// 创建合同草稿，并发起；由于在业务分类中配置好了签署动作、签署位置、合同文档等信息，所以此处无需额外指定
	contractId, err := draftByCategory(sdkClient)
	fmt.Println("创建合同草稿：", contractId, err)
	// 签署合同：本方企业签署合同
	err = signByCategory(sdkClient, contractId)
	fmt.Println("签署公章：", err)
	/**
	 * 本方公司签署完成，其他签署方签署可采用：
	 * （1）接收短信的方式登录契约锁云平台进行签署
	 * （2）生成内嵌页面签署链接进行签署（下方生成的链接）
	 */
	url, err := personalSignUrl(sdkClient, contractId)
	fmt.Println("个人签署链接：", url, err)
}

func draftByCategory(sdkClient *http.SdkClient) (contractId string, err error) {
	contract := model.Contract{}
	contract.Subject = "go测试合同-0821-1"
	user := model.User{}
	user.Name = "宋一"
	user.Contact = "10000000001"
	user.ContactType = "MOBILE"
	contract.Creator = &user // 发起人
	category := model.Category{}
	category.Name = "平台-个人-预设流程"
	contract.Category = &category // 业务分类
	// 指定模板参数：业务分类已绑定模板，此处只需要指定模板参数即可
	var templateParams []*model.TemplateParam
	param1 := model.TemplateParam{}
	param1.Name = "param1"
	param1.Value = "v1"
	param2 := model.TemplateParam{}
	param2.Name = "param2"
	param2.Value = "v2"
	templateParams = append(templateParams, &param1)
	templateParams = append(templateParams, &param2)
	contract.TemplateParams = templateParams
	// 设置合同签署方
	// 公司签署方--本方企业
	signatory1 := model.Signatory{}
	signatory1.TenantName = "测试11-8-1"
	signatory1.TenantType = "COMPANY"
	signatory1.Receiver = &user // 经办人
	signatory1SerialNo := 1
	signatory1.SerialNo = &signatory1SerialNo
	// 个人签署方
	signatory2 := model.Signatory{}
	signatory2.TenantName = "宋一"
	signatory2.TenantType = "PERSONAL"
	receiver := model.User{}
	receiver.Name = "宋一"
	receiver.Contact = "10000000001"
	receiver.ContactType = "MOBILE"
	signatory2.Receiver = &receiver
	signatory2SerialNo := 2
	signatory2.SerialNo = &signatory2SerialNo

	var signatories []*model.Signatory
	signatories = append(signatories, &signatory1)
	signatories = append(signatories, &signatory2)
	contract.Signatories = signatories
	send := true
	contract.Send = &send

	req := request.ContractDraftRequest{}
	req.Contract = &contract
	response, err := sdkClient.Service(req)
	if err != nil {
		return
	}
	code := response["code"].(float64)
	if code != 0 {
		msg := response["message"]
		err = errors.New(fmt.Sprintf("request fail,code:%.0f, message: %s", code, msg))
	} else {
		result := response["result"].(map[string]interface{})
		contractId = result["id"].(string)
	}
	return
}

func signByCategory(sdkClient *http.SdkClient, contractId string) (err error) {
	req := request.ContractSignCompanyRequest{}
	param := model.SignParam{}
	param.ContractId = contractId
	req.Param = &param

	response, err := sdkClient.Service(req)
	if err != nil {
		return
	}
	code := response["code"].(float64)
	if code != 0 {
		msg := response["message"]
		err = errors.New(fmt.Sprintf("request fail,code:%.0f, message: %s", code, msg))
	}
	return
}

func personalSignUrl(sdkClient *http.SdkClient, contractId string) (url string, err error) {
	req := request.ContractPageRequest{}
	req.ContractId = contractId
	user := model.User{}
	user.Contact = "10000000001"
	user.ContactType = "MOBILE"
	req.User = &user

	response, err := sdkClient.Service(req)
	if err != nil {
		return
	}
	code := response["code"].(float64)
	if code != 0 {
		msg := response["message"]
		err = errors.New(fmt.Sprintf("request fail,code:%.0f, message: %s", code, msg))
	} else {
		result := response["result"].(map[string]interface{})
		url = result["pageUrl"].(string)
	}
	return
}
