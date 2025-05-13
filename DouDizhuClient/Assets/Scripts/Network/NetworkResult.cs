namespace Network
{
    
    public struct NetworkResult
    {
        public bool IsSuccess { get; }
        public string ErrorMessage { get; }

        private NetworkResult(bool isSuccess, string errorMessage)
        {
            IsSuccess = isSuccess;
            ErrorMessage = errorMessage;
        }

        public static NetworkResult Success() => new NetworkResult(true, null);
        public static NetworkResult Failure(string errorMessage) => new NetworkResult(false, errorMessage);
    }

    public struct NetworkResult<T>
    {
        public bool IsSuccess { get; }
        public T Data { get; }
        public string ErrorMessage { get; }

        private NetworkResult(bool isSuccess, T data, string errorMessage)
        {
            IsSuccess = isSuccess;
            Data = data;
            ErrorMessage = errorMessage;
        }

        public static NetworkResult<T> Success(T data) => new NetworkResult<T>(true, data, null);
        public static NetworkResult<T> Failure(string errorMessage) => new NetworkResult<T>(false, default, errorMessage);
    }
}
