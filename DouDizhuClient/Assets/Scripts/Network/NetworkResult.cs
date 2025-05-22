namespace Network
{
    public struct NetworkResult<T>
    {
        public bool IsSuccess { get; set; }
        public T Data { get; set; }
        public string ErrorCode { get; set; }
        public string ErrorMessage { get; set; }

        private NetworkResult(bool isSuccess, T data, string errorCode, string errorMessage)
        {
            IsSuccess = isSuccess;
            Data = data;
            ErrorCode = errorCode;
            ErrorMessage = errorMessage;
        }

        public static NetworkResult<T> Success(T data) => new NetworkResult<T>(true, data, null, null);
        public static NetworkResult<T> Error(string errorCode, string errorMessage) => new NetworkResult<T>(false, default, errorCode, errorMessage);
    }
}