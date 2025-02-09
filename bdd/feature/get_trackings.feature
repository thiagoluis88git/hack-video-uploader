Feature: Upload video
  In order to get the list of my trackings
  As a user of the video uploader API
  I need to be able to see all my trackings

  Scenario: then user try to get the list, success list should be displayed
    When I send "GET" request to "/api/trackings/{cpf}" without payload:
      """  
      """
    Then the response code should be 200
    And the response payload should match json:
      """
        [
          {
            "trackingId": "d35863a5-85d6-4db7-8b06-075000a2948f",
            "status": "DONE",
            "videoUrl": "https://hack-video-processing-bucket.s3.amazonaws.com/d35863a5-85d6-4db7-8b06-075000a2948f",
            "zipUrl": "https://hack-processed-zip-bucket.s3.amazonaws.com/d35863a5-85d6-4db7-8b06-075000a2948f.zip",
            "createdAt": "2025-02-02T15:49:52.313546-03:00",
            "updatedAt": "2025-02-02T15:50:23.806347-03:00"
          },
          {
            "trackingId": "57db8998-2184-4e7a-871b-f70353a4cd61",
            "status": "DONE",
            "videoUrl": "https://hack-video-processing-bucket.s3.amazonaws.com/57db8998-2184-4e7a-871b-f70353a4cd61",
            "zipUrl": "https://hack-processed-zip-bucket.s3.amazonaws.com/57db8998-2184-4e7a-871b-f70353a4cd61.zip",
            "createdAt": "2025-02-02T15:51:43.789108-03:00",
            "updatedAt": "2025-02-02T15:52:07.026185-03:00"
          }
        ]
      """