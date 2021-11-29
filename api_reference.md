# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [proto/v1/dogfood/dogfood.proto](#proto/v1/dogfood/dogfood.proto)
    - [CreateRecordRequest](#dogfoodpb.v1.CreateRecordRequest)
    - [ListRecordsRequest](#dogfoodpb.v1.ListRecordsRequest)
    - [ListRecordsResponse](#dogfoodpb.v1.ListRecordsResponse)
    - [Record](#dogfoodpb.v1.Record)
  
    - [DogFoodService](#dogfoodpb.v1.DogFoodService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="proto/v1/dogfood/dogfood.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/v1/dogfood/dogfood.proto



<a name="dogfoodpb.v1.CreateRecordRequest"></a>

### CreateRecordRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dogfood_name | [string](#string) |  | dog_food name is a name of dogfood brand. |
| gram | [int32](#int32) |  | grap specifies how grams a dog eat dogfood. |
| dog_name | [string](#string) |  | dog_name specifies a name of dog. |






<a name="dogfoodpb.v1.ListRecordsRequest"></a>

### ListRecordsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| from | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | from specifies the start time of eaten_at. |
| page_size | [int32](#int32) |  | page_size specifies a requested length of records. |
| to | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | to specifies the end time of eaten_at. |






<a name="dogfoodpb.v1.ListRecordsResponse"></a>

### ListRecordsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| records | [Record](#dogfoodpb.v1.Record) | repeated | records specify an array of Record. |
| to | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | to specifies the end time of eaten_at. |






<a name="dogfoodpb.v1.Record"></a>

### Record



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dogfood_name | [string](#string) |  | dog_food name is a name of dogfood brand. |
| gram | [int32](#int32) |  | grap specifies how grams a dog eat dogfood. |
| dog_name | [string](#string) |  | dog_name specifies a name of dog. |
| eaten_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | eaten_at specifies what time a dog ate a dogfood. |





 

 

 


<a name="dogfoodpb.v1.DogFoodService"></a>

### DogFoodService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateRecord | [CreateRecordRequest](#dogfoodpb.v1.CreateRecordRequest) | [Record](#dogfoodpb.v1.Record) | CreateRecord create a record who ate what, when, and how much. |
| ListRecords | [ListRecordsRequest](#dogfoodpb.v1.ListRecordsRequest) | [ListRecordsResponse](#dogfoodpb.v1.ListRecordsResponse) | ListRecords list up all records. |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

