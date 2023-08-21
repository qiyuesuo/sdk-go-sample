package sample

import (
	"errors"
	"fmt"
	"os"
	"qiyuesuo/sdk/http"
	"qiyuesuo/sdk/model"
	"qiyuesuo/sdk/request"
	"testing"
)

/*
签署示例：该场景模拟一个人事合同的场景，即平台方公司与员工签署合同，平台方公司先签署，员工后签。
具体操作：
1、自定义签署流程，指定好签署方、签署流程、签署位置、合同文档等信息，并发起合同；
2、通过接口签署本方企业，再通过接口获取第三方个人的签署链接，个人打开签署链接签署合同。
*/
func TestCustomSend(t *testing.T) {
	sdkClient := http.NewSdkClient("https://openapi.qiyuesuo.cn", "tBYw111111", "yHnVOk11MNUT97BKQnN5xnul111111")
	//sdkClient := http.NewSdkClient("更换为开放平台请求地址", "更换为您开放平台 App Token", "更换为您开放平台App Secret")
	// 创建合同草稿：指定签署方、签署动作等信息
	contractId, contractDetail, err := customDraft(sdkClient)
	fmt.Println("创建合同草稿：", contractId, err)
	// 添加合同文档
	documentId, err := addDocumentByFile(sdkClient, contractId)
	fmt.Println("添加合同文档：", documentId, err)
	// 添加合同文档
	documentId2, err := addDocumentByTemplate(sdkClient, contractId)
	fmt.Println("用模板添加合同文档：", documentId2, err)
	// 发起合同，同时指定签署位置待信息
	err = send(sdkClient, contractId, contractDetail, documentId)
	fmt.Println("发起合同：", err)
	// 签署合同：本方企业签署合同
	err = customSign(sdkClient, contractId)
	fmt.Println("签署公章：", err)
	/**
	 * 本方公司签署完成，其他签署方签署可采用：
	 * （1）接收短信的方式登录契约锁云平台进行签署
	 * （2）生成内嵌页面签署链接进行签署（下方生成的链接）
	 */
	url, err := generatePersonalSignUrl(sdkClient, contractId)
	fmt.Println("个人签署链接：", url, err)
}

func customDraft(sdkClient *http.SdkClient) (contractId string, contractDetail map[string]interface{}, err error) {
	contract := model.Contract{}
	contract.Subject = "go测试合同-0818-1"
	user := model.User{}
	user.Name = "宋一"
	user.Contact = "10000000001"
	user.ContactType = "MOBILE"
	contract.Creator = &user // 发起人
	// 设置合同签署方
	// 公司签署方--本方企业
	signatory1 := model.Signatory{}
	signatory1.TenantName = "测试11-8-1"
	signatory1.TenantType = "COMPANY"
	signatory1.Receiver = &user // 经办人
	signatory1SerialNo := 1
	signatory1.SerialNo = &signatory1SerialNo
	// 公章签署流程
	sealAction := model.Action{}
	actionSn := 1
	sealAction.SerialNo = &actionSn
	sealAction.Type_ = "COMPANY"
	sealAction.SealId = "3124619033605980940"
	var sig1Actions []*model.Action
	sig1Actions = append(sig1Actions, &sealAction)
	signatory1.Actions = sig1Actions

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
	send := false
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
		contractDetail = result
		contractId = contractDetail["id"].(string)
	}
	return
}

func addDocumentByFile(sdkClient *http.SdkClient, contractId string) (documentId string, err error) {
	req := request.DocumentAddByFileRequest{}
	req.ContractId = contractId
	req.Title = "goFile1"
	req.FileSuffix = "pdf"
	file, _ := os.Open("/Users/sss/Downloads/文件模板三.pdf")
	req.File = &http.FileItem{file, ""}

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
		documentId = result["documentId"].(string)
	}
	return
}

func addDocumentByTemplate(sdkClient *http.SdkClient, contractId string) (documentId string, err error) {

	req := request.DocumentAddByTemplateRequest{}
	req.ContractId = contractId
	req.Title = "goTemplateFile"
	req.TemplateId = "2766933842559242287"
	var templateParams []*model.TemplateParam
	param1 := model.TemplateParam{}
	param1.Name = "param1"
	param1.Value = "v1"
	param2 := model.TemplateParam{}
	param2.Name = "param2"
	param2.Value = "v2"
	templateParams = append(templateParams, &param1)
	templateParams = append(templateParams, &param2)
	req.TemplateParams = templateParams

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
		documentId = result["documentId"].(string)
	}
	return
}

func send(sdkClient *http.SdkClient, contractId string, contractDetail map[string]interface{}, documentId string) (err error) {
	// 获取SignatoryId与ActionId，用于指定签署位置，公司签署位置需要指定ActionId,个人签署位置需要指定SignatoryId
	var personalSignatoryId, companySealActionId string
	signatories := contractDetail["signatories"].([]interface{})
	for _, v := range signatories {
		sig := v.(map[string]interface{})
		tenantType := sig["tenantType"].(string)
		if tenantType == "PERSONAL" {
			// 取个人签署方 ID
			personalSignatoryId = sig["id"].(string)
		} else {
			// 从企业签署方中取公章签署动作 ID
			companyActions := sig["actions"].([]interface{})
			for _, caV := range companyActions {
				ca := caV.(map[string]interface{})
				caType := ca["type"].(string)
				if caType == "COMPANY" {
					companySealActionId = ca["id"].(string)
				}
			}
		}
	}

	req := request.ContractSendRequest{}
	req.ContractId = contractId
	// 签署位置-公章签署
	companyStamper := model.Stamper{}
	companyStamper.Type_ = "COMPANY"
	companyStamper.ActionId = companySealActionId
	companyStamper.DocumentId = documentId
	page := 1
	companyStamper.Page = &page
	offsetX := 0.1
	companyStamper.OffsetX = &offsetX
	offsetY := 0.1
	companyStamper.OffsetY = &offsetY
	// 签署位置-个人签署
	personalStamper := model.Stamper{}
	personalStamper.Type_ = "PERSONAL"
	personalStamper.SignatoryId = personalSignatoryId
	personalStamper.DocumentId = documentId
	pPage := 1
	personalStamper.Page = &pPage
	pOffsetX := 0.3
	personalStamper.OffsetX = &pOffsetX
	pOffsetY := 0.1
	personalStamper.OffsetY = &pOffsetY
	var stampers []*model.Stamper
	stampers = append(stampers, &companyStamper)
	stampers = append(stampers, &personalStamper)
	req.Stampers = stampers
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

func customSign(sdkClient *http.SdkClient, contractId string) (err error) {
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

func generatePersonalSignUrl(sdkClient *http.SdkClient, contractId string) (url string, err error) {
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
