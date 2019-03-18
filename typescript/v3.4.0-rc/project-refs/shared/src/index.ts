export { hello } from "../../core/dist";

export function addSuffix<T extends (...args: any[]) => string>(f: T, suffix: string): T {
    return ((...args: any[]) => `${f(...args)}${suffix}`) as T;
}
