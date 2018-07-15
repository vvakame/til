// defaultPropsを考慮したPropsを計算する
type Defaultize<TProps, TDefaults> =
    // defaultPropsに含まれるプロパティはoptionalにする
    & { [K in Extract<keyof TProps, keyof TDefaults>]?: TProps[K] }
    // defaultPropsに含まれないものはそのまま
    & { [K in Exclude<keyof TProps, keyof TDefaults>]: TProps[K] }
    // defaultPropsにしか含まれないものはoptionalにする
    & Partial<TDefaults>;

// propTypesに指定されたPropTypesの値から求められる型を計算する
type InferredPropTypes<P> = { [K in keyof P]: P[K] extends PropTypeChecker<infer T, infer U> ? PropTypeChecker<T, U>[typeof checkedType] : {} };

// 型計算用のシンボルを定義
declare const checkedType: unique symbol;
// propTypesに定義する値 何の型を要求するかとそれがrequiredかどうか
interface PropTypeChecker<U, TRequired = false> {
    (props: any, propName: string, componentName: string, location: any, propFullName: string): boolean;
    isRequired: PropTypeChecker<U, true>;
    [checkedType]: TRequired extends true ? U : U | null | undefined;
}

// 値空間に存在するPropTypesの定義
declare namespace PropTypes {
    export const number: PropTypeChecker<number>;
    export const string: PropTypeChecker<string>;
    export const node: PropTypeChecker<ReactNode>;
}

type ReactNode = string | number | ReactComponent<{}, {}>;

declare class ReactComponent<P={}, S={}> {
    constructor(props: P);
    props: P & Readonly<{ children: ReactNode[] }>;
    setState(s: Partial<S>): S;
    render(): ReactNode;
}

declare namespace JSX {
    interface Element extends ReactComponent { }
    interface IntrinsicElements { }

    type LibraryManagedAttributes<TComponent, TProps> =
        // defaultProps, propTypes の両方が存在するか
        TComponent extends { defaultProps: infer D; propTypes: infer P; }
            ? Defaultize<TProps & InferredPropTypes<P>, D>
            // defaultProps が存在するか
            : TComponent extends { defaultProps: infer D }
                ? Defaultize<TProps, D>
                // propTypes が存在するか
                : TComponent extends { propTypes: infer P }
                    ? TProps & InferredPropTypes<P>
                    : TProps;
}
