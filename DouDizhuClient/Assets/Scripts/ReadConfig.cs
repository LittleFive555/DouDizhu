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
        var value = ConfigsManager.Instance.GetConst<int>("AccountMinLength");
        Debug.Log(value);
    }

    // Update is called once per frame
    void Update()
    {
        
    }
}
