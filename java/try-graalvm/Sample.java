// package net.vvakame.graalvm;

public class Sample {
    public static void main(String[] args) {
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < 10000000; i++) {
            sb.append("test");
        }

        System.out.println(sb.toString().length());
    }
}
