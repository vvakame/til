# マジでちょっとだけGraalVMを触る

https://www.graalvm.org/downloads/
Community EditionとEnterprise Editionの両方あるのがOracleらしい…。

> Free for development and production use

らしいのでCommunity Edition使う。

https://hub.docker.com/r/oracle/graalvm-ce
`oracle/graalvm-ce:latest` なるほど。
今は `docker pull oracle/graalvm-ce:19.2.0.1` と等しいのかな。

すげー、Docker Hubにライセンスについてたくさん書いてあって使い方が一ミリも書いてない。

https://www.graalvm.org/docs/examples/native-list-dir/ とか見たほうがよさそう。

```
$ cat Sample.java
public class Sample {
    public static void main(String[] args) {
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < 10000000; i++) {
            sb.append("test");
        }

        System.out.println(sb.toString().length());
    }
}

$ docker run -it --rm -v $(pwd):/tmp/javacode oracle/graalvm-ce:latest /bin/bash
# gu install native-image
# cd /tmp/javacode
# javac Sample.java
# native-image Sample
Build on Server(pid: 114, port: 43351)
[sample:114]    classlist:     787.38 ms
[sample:114]        (cap):   1,374.73 ms
[sample:114]        setup:   3,223.41 ms
[sample:114]   (typeflow):   3,275.53 ms
[sample:114]    (objects):   1,982.55 ms
[sample:114]   (features):     209.05 ms
[sample:114]     analysis:   5,568.89 ms
[sample:114]     (clinit):     121.96 ms
[sample:114]     universe:     493.35 ms
[sample:114]      (parse):     528.75 ms
[sample:114]     (inline):   1,146.52 ms
[sample:114]    (compile):   5,354.03 ms
[sample:114]      compile:   7,338.53 ms
[sample:114]        image:     449.87 ms
[sample:114]        write:     155.12 ms
[sample:114]      [total]:  18,256.71 ms

# time ./sample
40000000

real	0m0.423s
user	0m0.310s
sys	0m0.090s

# time java Sample
40000000

real	0m0.533s
user	0m0.500s
sys	0m0.270s

# file ./sample
./sample: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked (uses shared libs), for GNU/Linux 2.6.32, BuildID[sha1]=416334bb53a39c007a4e5c074cf7f2d50c3de5da, not stripped

# ldd ./sample
	linux-vdso.so.1 =>  (0x00007fff425d9000)
	libm.so.6 => /lib64/libm.so.6 (0x00007f6192ba7000)
	libdl.so.2 => /lib64/libdl.so.2 (0x00007f61929a3000)
	libpthread.so.0 => /lib64/libpthread.so.0 (0x00007f6192787000)
	libz.so.1 => /lib64/libz.so.1 (0x00007f6192571000)
	librt.so.1 => /lib64/librt.so.1 (0x00007f6192369000)
	libc.so.6 => /lib64/libc.so.6 (0x00007f6191f9b000)
	/lib64/ld-linux-x86-64.so.2 (0x00007f6192ea9000)

# native-image --static Sample

# file ./sample
./sample: ELF 64-bit LSB executable, x86-64, version 1 (GNU/Linux), statically linked, for GNU/Linux 2.6.32, BuildID[sha1]=7daddca3016819b0db903a28ac1bbb15e8dd4d26, not stripped
```

ふーんぬ。
`--enable-http` とか `--enable-https` とか色々あるなー。
面白いけど今のところ若干沼を感じる。
マジでチューニングし始めたら2-3日くらいマニュアルやらリファレンスやら読みまくって設定を弄くり続ける必要を感じる。
RubyとかNode.jsとかも対応していて、あんまりNative binary化は推しポイントではないっぽい感じもする。

あんましJVMが必要になるタイミングがないんだけど必要になったら真面目に検討してもよさそう。
