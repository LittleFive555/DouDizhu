using System;
using Org.BouncyCastle.Asn1.Sec;
using Org.BouncyCastle.Crypto;
using Org.BouncyCastle.Crypto.Agreement;
using Org.BouncyCastle.Crypto.Digests;
using Org.BouncyCastle.Crypto.Engines;
using Org.BouncyCastle.Crypto.Generators;
using Org.BouncyCastle.Crypto.Modes;
using Org.BouncyCastle.Crypto.Paddings;
using Org.BouncyCastle.Crypto.Parameters;
using Org.BouncyCastle.Math;
using Org.BouncyCastle.Security;

namespace Network.Encryption
{
    /// <summary>
    /// 加密器
    /// <para>使用ECDH算法生成共享密钥，使用HKDF算法派生安全密钥</para>
    /// <para>使用AES-256-CBC算法加密和解密数据</para>
    /// </summary>
    public class Encryptor
    {
        public bool IsGenerated => m_DerivedSecureKey != null;

        private ECPrivateKeyParameters privateKey;
        private ECDomainParameters domainParams;
        private byte[] m_DerivedSecureKey;
        
        public byte[] GeneratePublicKey()
        {

            // 1. 生成客户端临时密钥对
            var ecParams = SecNamedCurves.GetByName("secp256r1");
            domainParams = new ECDomainParameters(
                ecParams.Curve, ecParams.G, ecParams.N, ecParams.H);
            
            var keyPair = GenerateKeyPair(domainParams);
            privateKey = (ECPrivateKeyParameters)keyPair.Private;
            var publicKey = (ECPublicKeyParameters)keyPair.Public;
            return publicKey.Q.GetEncoded(false);
        }

        public void DeriveSecureKey(byte[] serverPublicKeyBytes, byte[] salt, byte[] info)
        {
            var sharedSecret = DeriveSharedSecret(serverPublicKeyBytes);
            m_DerivedSecureKey = DeriveSecureKeyImpl(sharedSecret, salt, info, 32);
        }

        public (byte[] iv, byte[] ciphertext) Encrypt(byte[] plaintext)
        {
            // 生成随机IV(块大小为16字节)
            byte[] iv = new byte[16];
            new SecureRandom().NextBytes(iv);

            // 初始化加密引擎 CBC模式
            AesEngine engine = new AesEngine();  // 不建议存储引擎实例
            PaddedBufferedBlockCipher cipher = new PaddedBufferedBlockCipher(new CbcBlockCipher(engine), new Pkcs7Padding());
            cipher.Init(true, new ParametersWithIV(new KeyParameter(m_DerivedSecureKey), iv));

            // 执行加密
            byte[] output = new byte[cipher.GetOutputSize(plaintext.Length)];
            int len = cipher.ProcessBytes(plaintext, 0, plaintext.Length, output, 0);
            cipher.DoFinal(output, len);

            return (iv, output);
        }

        public byte[] Decrypt(byte[] ciphertext, byte[] iv)
        {
            AesEngine engine = new AesEngine();
            PaddedBufferedBlockCipher cipher = new PaddedBufferedBlockCipher(new CbcBlockCipher(engine), new Pkcs7Padding());
            cipher.Init(false, new ParametersWithIV(new KeyParameter(m_DerivedSecureKey), iv));

            byte[] output = new byte[cipher.GetOutputSize(ciphertext.Length)];
            int len = cipher.ProcessBytes(ciphertext, 0, ciphertext.Length, output, 0);
            int finalLen = cipher.DoFinal(output, len);
            
            // 创建正确大小的数组，只包含实际的明文数据（去除填充）
            byte[] result = new byte[len + finalLen];
            Array.Copy(output, 0, result, 0, len + finalLen);
            
            return result;
        }

        private AsymmetricCipherKeyPair GenerateKeyPair(ECDomainParameters domainParams)
        {
            var generator = new ECKeyPairGenerator();
            var keyGenParams = new ECKeyGenerationParameters(domainParams, new SecureRandom());
            generator.Init(keyGenParams);
            return generator.GenerateKeyPair();
        }

        private byte[] DeriveSharedSecret(byte[] otherPartyPublicKeyBytes)
        {
            // 导入对方公钥
            var curve = domainParams.Curve;
            var otherPartyPoint = curve.DecodePoint(otherPartyPublicKeyBytes);
            var otherPartyPublicKey = new ECPublicKeyParameters(otherPartyPoint, domainParams);
            
            // 计算共享密钥
            var agreement = new ECDHBasicAgreement();
            agreement.Init(privateKey);
            var sharedSecret = agreement.CalculateAgreement(otherPartyPublicKey);
            
            return ToFixedLengthBytes(sharedSecret, 32); // 32 字节对齐
        }
        
        byte[] ToFixedLengthBytes(BigInteger value, int length)
        {
            byte[] bytes = value.ToByteArrayUnsigned();
            if (bytes.Length == length) return bytes;
            
            byte[] result = new byte[length];
            Buffer.BlockCopy(
                src: bytes,
                srcOffset: Math.Max(0, bytes.Length - length),
                dst: result,
                dstOffset: Math.Max(0, length - bytes.Length),
                count: Math.Min(bytes.Length, length)
            );
            return result;
        }

        private byte[] DeriveSecureKeyImpl(byte[] sharedSecret, byte[] salt, byte[] info, int outputLength)
        {
            HkdfParameters parameters = new HkdfParameters(sharedSecret, salt, info);
            HkdfBytesGenerator hkdf = new HkdfBytesGenerator(new Sha256Digest());
            hkdf.Init(parameters);
            byte[] output = new byte[outputLength];
            hkdf.GenerateBytes(output, 0, outputLength);
            return output;
        }
    }
}