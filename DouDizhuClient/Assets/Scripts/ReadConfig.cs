using UnityEngine;
using Config;
using Config.Define;

public class ReadConfig : MonoBehaviour
{
    private void Awake()
    {
        ConfigsManager.Instance.LoadConfigs();
    }
    // Start is called before the first frame update
    void Start()
    {
        var config = ConfigsManager.Instance.GetConfig<DConst>("AccountMinLength");
        Debug.Log(config.Value + " " + config.Type);
    }

    // Update is called once per frame
    void Update()
    {
        
    }
}
