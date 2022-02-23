import json
from configuration.env_config import LOG_PATH, INPUT_BASE, orchestrator, clm, dcm, dtc, hpa, ncm, HOST



class InputParser(object):
    def __init__(self, file_name):
        self.file_name = file_name
        with open(file_name, 'r') as data_file:
            try:
                self.data = json.load(data_file)
            except Exception as e:
                print("\n================================================================\n")
                print("[ERROR]: Info:Could not load testcases from file '{0}'; Exception: '{1}' ".format(file_name, str(e)))
                print("\n================================================================\n")

    def get_data_size(self):
        try:
            return len(self.data)
        except KeyError:
            print("Invalid data  " + self.data)
            return None

    def get_type(self, index):
        try:
            return self.data[index]['type']
        except KeyError:
            print("Info:Could not find key type")
            return None

    def get_type_metadata(self, index):
        try:
            return self.data[index]['type-metadata']
        except KeyError:
            print("Info:Could not find key type-metadata")
            return None
    
    def get_type_metadata_length(self, index):
        try:
            return len(self.data[index]['type-metadata'])
        except KeyError:
            print("Info:Could not find length of key type-metadata")
            return None

    def get_request_body(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['RequestBody']
        except KeyError:
            print("Info:Could not find key RequestBody")
            return None
    
    def get_request_body_controller(self, index, sub_index):
        try:
            reqBody = self.data[index]['type-metadata'][sub_index]['RequestBody']
            # with open('MASTER_NODE_IP.txt', 'r') as f:
            #     host_ip = f.read().strip("\n")
            host_ip = HOST
            reqBody["spec"]["host"] = host_ip
            return reqBody

        except KeyError:
            print("Info:Could not find key RequestBody")
            return None

    def get_request_body_name(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['RequestBody']['metadata']['name']
        except KeyError:
            print("Info:Could not find key RequestBody metadata 'name'")
            return None

    def get_request_body_version(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['RequestBody']['spec']['version']
        except KeyError:
            print("Info:Could not find key RequestBody metadata 'name'")
            return None
    
    def get_request_body_labelname(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['RequestBody']['label-name']
        except KeyError:
            print("Info:Could not find key RequestBody 'label-name'")
            return None

    def get_anchor(self, index):
        try:
            return self.data[index]['anchor']
        except KeyError:
            print("Info:Could not find anchor")
            return None

    def get_request_body_file(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['RequestBodyFile']
        except KeyError:
            print("Info:Could not find key file")
            return None
    
    def get_response_code(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ResponseCode']
        except KeyError:
            print("Info:Could not find key ResponseCode")
            return None

    def get_response_text(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ResponseText']
        except KeyError:
            print("Info:Could not find key ResponseText")
            return None
    
    def get_response_text_controller(self, index, sub_index):
        try:
            resText = self.data[index]['type-metadata'][sub_index]['ResponseText']
            # with open('MASTER_NODE_IP.txt', 'r') as f:
            #     host_ip = f.read().strip("\n")
            host_ip = HOST
            resText["spec"]["host"] = host_ip
            return resText

        except KeyError:
            print("Info:Could not find key ResponseText")
            return None

    # def get_response_text_controller_getAll(self, index, sub_index):
    #     try:
    #         resText = self.data[index]['type-metadata'][sub_index]['ResponseText']
    #         with open('MASTER_NODE_IP.txt', 'r') as f:
    #             host_ip = f.read().strip("\n")
    #         resText[0]["spec"]["host"] = host_ip
    #         return resText

    #     except KeyError:
    #         print("Info:Could not find key ResponseText")
    #         return None

    def get_response_text_controller_getAll(self, index, sub_index):
        try:
            resText = self.data[index]['type-metadata'][sub_index]['ResponseText']
            # with open('MASTER_NODE_IP.txt', 'r') as f:
            #     host_ip = f.read().strip("\n")
            host_ip = HOST
            if len(resText) > 1:
                for response in resText:
                    response["spec"]["host"] = host_ip
            else:
                resText[0]["spec"]["host"] = host_ip
            print(resText)
            return resText
        except KeyError:
            print("Info:Could not find key ResponseText")
            return None
    
    def get_response_text_for_GET(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ResponseTextForGET']
        except KeyError:
            print("Info:Could not find key ResponseText")
            return None

    def get_param_projectName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamProjectName']
        except KeyError:
            print("Info:Could not find key ParamProjectName")
            return None

    def get_param_controllerName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamControllerName']
        except KeyError:
            print("Info:Could not find key ParamControllerName")
            return None
    
    def get_param_clusterProviderName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamClusterProviderName']
        except KeyError:
            print("Info:Could not find key ParamClusterProviderName")
            return None
 
    def get_param_clusterName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamClusterName']
        except KeyError:
            print("Info:Could not find key ParamClusterName")
            return None

    def get_param_clusterLabelName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamClusterLabelName']
        except KeyError:
            print("Info:Could not find key ParamClusterLabelName")
            return None

    def get_param_clusterKVPair(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamClusterKVPair']
        except KeyError:
            print("Info:Could not find key ParamClusterKVPair")
            return None

    def get_param_compositeAppName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamCompositeAppName']
        except KeyError:
            print("Info:Could not find key ParamCompositeAppName")
            return None

    def get_param_compositeAppVersion(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamCompositeAppVersion']
        except KeyError:
            print("Info:Could not find key ParamCompositeAppVersion")
            return None

    def get_param_compositeProfileName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamCompositeProfileName']
        except KeyError:
            print("Info:Could not find key ParamCompositeProfileName")
            return None

    def get_param_deploymentIntentGroupName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamDeploymentIntentGroupName']
        except KeyError:
            print("Info:Could not find key ParamDeploymentIntentGroupName")
            return None

    def get_param_intentName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamIntentName']
        except KeyError:
            print("Info:Could not find key ParamIntentName")
            return None

    def get_param_genericPlacementIntent(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamGenericPlacementIntent']
        except KeyError:
            print("Info:Could not find key ParamGenericPlacementIntent")
            return None

    def get_param_intentGenericPlacementIntent(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamIntentGenericPlacementIntent']
        except KeyError:
            print("Info:Could not find key ParamIntentGenericPlacementIntent")
            return None

    def get_param_appName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamAppName']
        except KeyError:
            print("Info:Could not find key ParamAppName")
            return None
    
    def get_param_networkName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamNetworkName']
        except KeyError:
            print("Info:Could not find key ParamNetworkName")
            return None
    
    def get_param_appProfileName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamAppProfileName']
        except KeyError:
            print("Info:Could not find key ParamAppProfileName")
            return None
    
    def get_param_logicalCloudName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamLogicalCloudName']
        except KeyError:
            print("Info:Could not find key ParamLogicalCloudName")
            return None

    def get_param_clusterReferenceName(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamClusterReferenceName']
        except KeyError:
            print("Info:Could not find key ParamClusterReferenceName")
            return None
    
    def get_param_getStatusOptionalParams(self, index, sub_index):
        try:
            return self.data[index]['type-metadata'][sub_index]['ParamOptionalGetStatus']
        except KeyError:
            print("Info:Could not find key ParamOptionalGetStatus")
            return None


    '''
    @staticmethod
    def get_no_of_subapps(response):
        return len(response['subapps'])

    @staticmethod
    def get_subapp_index(response, subapp_name):
        for _subapps in response['subapps']:
            for key, value in _subapps.items():
                if value == subapp_name:
                    return(response['subapps'].index(_subapps))
    '''
    
    

    def get_no_of_post_project_items(self):
        return len(self.data[0]['post-projects'])

    def get_no_of_get_project_items(self):
        return len(self.data[1]['get-projects'])



