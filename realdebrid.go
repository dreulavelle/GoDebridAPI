package GoDebridAPI

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

var api_url = url.URL{
	Scheme: "https",
	Host:   "api.real-debrid.com",
	Path:   "/rest/",
}

type Client struct {
	c      *http.Client
	apiKey string
}

/* endpoint /user
* Params:
* returns: All the details about the user in Json
 */
func (c *Client) RdGetUser() (rdUserSchema, error) {
	var user rdUserSchema
	var wg sync.WaitGroup
	wg.Add(1)

	var err error
	go func() {
		defer wg.Done()
		err = c.CommonGetReq("/user", &user)
	}()
	wg.Wait()

	if err != nil {
		return rdUserSchema{}, err
	}
	return user, nil
}

/* endpoint /torrents
* Params:
* returns: all the torrents added by the user
 */
func (c *Client) rdGetTorrents() ([]rdTorrentSchema, error) {
	resBody, err := c.GetReq("/torrents")
	if err != nil {
		return nil, fmt.Errorf("Couldnt get user")
	}

	torrents := []rdTorrentSchema{}
	if err := json.Unmarshal(resBody, &torrents); err != nil {
		return nil, fmt.Errorf("Decode failed")
	}

	return torrents, nil
}

/* endpoint /torrents/addMagnet
* Params: magnet link string
* returns: id and url of the torrent added
 */
func (c *Client) rdAddMagnet(magnet string) (rdAddMagnetSchema, error) {
	data := url.Values{}
	data.Set("magnet", magnet)

	resBody, err := c.PostReq("/torrents/addMagnet", data)
	if err != nil {
		return rdAddMagnetSchema{}, err
	}

	mag := rdAddMagnetSchema{}
	if err := json.Unmarshal(resBody, &mag); err != nil {
		return rdAddMagnetSchema{}, fmt.Errorf("Decode failed")
	}

	return mag, nil
}

/* endpoint /torrents/info/{id}.
* Params: Id of the torrrent whose info is needed
* returns: all the details of the torrent in json format
 */
func (c *Client) rdGetFileInfo(id string) (rdTorrentInfoSchema, error) {

	resBody, err := c.GetReq(fmt.Sprintf("/torrents/info/%s", id))
	if err != nil {
		return rdTorrentInfoSchema{}, err
	}

	fileInfo := rdTorrentInfoSchema{}
	if err := json.Unmarshal(resBody, &fileInfo); err != nil {
		return rdTorrentInfoSchema{}, fmt.Errorf("Decode failed")
	}

	return fileInfo, nil
}

/* endpoint /torrents/selectFiles/{id}.
* Params: Id of the torrent, we can get id from /torrents/info
* returns: Nothing
 */
func (c *Client) rdSelectFiles(id string) error {

	torrentFiles, err := c.rdGetFileInfo(id)
	if err != nil {
		fmt.Errorf("Couldnt get files from the torrent")
	}

	files := GetFileIdsFromTorrent(torrentFiles)

	data := url.Values{}
	data.Set("files", files)
	req, err := c.PostReq("/torrents/selectFiles/"+id, data)

	fmt.Println(string(req))

	if err != nil {
		fmt.Errorf("Couldnt make the request")
	}
	return nil
}

/* endpoint /unrestrict/link.
* Params: Link to the torrent
* returns: Unrestricted real-debrid link which can be downloaded by aria
 */
func (c *Client) rdUnrestrictLinks(link string) (UnrestrictLink, error) {
	data := url.Values{}
	data.Set("link", link)

	resp, err := c.PostReq("/unrestrict/link", data)
	if err != nil {
		return UnrestrictLink{}, err
	}

	getLink := UnrestrictLink{}
	if err := json.Unmarshal(resp, &getLink); err != nil {
		return UnrestrictLink{}, fmt.Errorf("Decode failed")
	}

	return getLink, nil
}

/* endpoint /downloads
* Params: None
* Returns: List of user's downloads
 */
func (c *Client) rdGetDownloads() ([]rdDownloadSchema, error) {
	resBody, err := c.GetReq("/downloads")
	if err != nil {
		return nil, fmt.Errorf("Couldn't get downloads")
	}

	downloads := []rdDownloadSchema{}
	if err := json.Unmarshal(resBody, &downloads); err != nil {
		return nil, fmt.Errorf("Decode failed")
	}

	return downloads, nil
}

/* endpoint /hosts
* Params: None
* Returns: List of supported hosts
 */
func (c *Client) rdGetHosts() ([]rdHostSchema, error) {
	resBody, err := c.GetReq("/hosts")
	if err != nil {
		return nil, fmt.Errorf("Couldn't get hosts")
	}

	hosts := []rdHostSchema{}
	if err := json.Unmarshal(resBody, &hosts); err != nil {
		return nil, fmt.Errorf("Decode failed")
	}

	return hosts, nil
}

/* endpoint /traffic
* Params: None
* Returns: User's traffic details
 */
func (c *Client) rdGetTraffic() (rdTrafficSchema, error) {
	resBody, err := c.GetReq("/traffic")
	if err != nil {
		return rdTrafficSchema{}, fmt.Errorf("Couldn't get traffic details")
	}

	traffic := rdTrafficSchema{}
	if err := json.Unmarshal(resBody, &traffic); err != nil {
		return rdTrafficSchema{}, fmt.Errorf("Decode failed")
	}

	return traffic, nil
}

/* endpoint /device/code
* Params: None
* Returns: Device code for user authentication
 */
func (c *Client) rdGetDeviceCode() (rdDeviceCodeSchema, error) {
	resBody, err := c.GetReq("/device/code")
	if err != nil {
		return rdDeviceCodeSchema{}, fmt.Errorf("Couldn't get device code")
	}

	deviceCode := rdDeviceCodeSchema{}
	if err := json.Unmarshal(resBody, &deviceCode); err != nil {
		return rdDeviceCodeSchema{}, fmt.Errorf("Decode failed")
	}

	return deviceCode, nil
}

/* endpoint /device/credentials
* Params: Device code
* Returns: Device credentials for user authentication
 */
func (c *Client) rdGetDeviceCredentials(deviceCode string) (rdDeviceCredentialsSchema, error) {
	data := url.Values{}
	data.Set("code", deviceCode)

	resBody, err := c.PostReq("/device/credentials", data)
	if err != nil {
		return rdDeviceCredentialsSchema{}, err
	}

	credentials := rdDeviceCredentialsSchema{}
	if err := json.Unmarshal(resBody, &credentials); err != nil {
		return rdDeviceCredentialsSchema{}, fmt.Errorf("Decode failed")
	}

	return credentials, nil
}

func (c *Client) rdGetTranscode(fileId string) (rdTranscodeSchema, error) {
	data := url.Values{}
	data.Set("id", fileId)

	resBody, err := c.PostReq("/streaming/transcode", data)
	if err != nil {
		return rdTranscodeSchema{}, err
	}

	transcode := rdTranscodeSchema{}
	if err := json.Unmarshal(resBody, &transcode); err != nil {
		return rdTranscodeSchema{}, fmt.Errorf("Decode failed")
	}

	return transcode, nil
}

/* endpoint /downloads/delete
* Params: Download ID
* Returns: Delete download result
 */
func (c *Client) rdDeleteDownload(downloadId string) (rdDeleteDownloadSchema, error) {
	resBody, err := c.GetReq("/downloads/delete/" + downloadId)
	if err != nil {
		return rdDeleteDownloadSchema{}, fmt.Errorf("Couldn't delete download")
	}

	deleteResult := rdDeleteDownloadSchema{}
	if err := json.Unmarshal(resBody, &deleteResult); err != nil {
		return rdDeleteDownloadSchema{}, fmt.Errorf("Decode failed")
	}

	return deleteResult, nil
}

/* endpoint /downloads/clear
* Params: None
* Returns: Clear downloads result
 */
func (c *Client) rdClearDownloads() (rdClearDownloadSchema, error) {
	resBody, err := c.GetReq("/downloads/clear")
	if err != nil {
		return rdClearDownloadSchema{}, fmt.Errorf("Couldn't clear downloads")
	}

	clearResult := rdClearDownloadSchema{}
	if err := json.Unmarshal(resBody, &clearResult); err != nil {
		return rdClearDownloadSchema{}, fmt.Errorf("Decode failed")
	}

	return clearResult, nil
}
