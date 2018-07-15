class Component extends ReactComponent {
    static propTypes = {
        foo: PropTypes.number,
        bar: PropTypes.node,
        baz: PropTypes.string.isRequired,
    };
    static defaultProps = {
        foo: 42,
    }
}

// ざっくりこういうイメージ 型なので = 42 とかは本来は書けない
// type Props = {
//     foo: number = 42;
//     bar: ReactNode;
//     baz: string;
// };

// OK
const a = <Component foo={12} bar="yes" baz="yeah" />;
const b = <Component bar="yes" baz="yeah" />;
const c = <Component foo={12} bar={null} baz="cool" />;

// エラーになる奴ら
// barはundefinedを受け付けるけどプロパティ自体を省略できるわけではない(≒型定義の不備に近い)
// const d = <Component foo={12} baz="yeah" />;
// bazがないのでエラー
// const e = <Component foo={12} bar={void 0} />;
// batはPropsの定義に存在しないのでエラー
// const f = <Component bar="yes" baz="yo" bat="ohno" />;
// bazはnon-nullなのでエラー
// const g = <Component foo={12} bar="yeah" baz={null} />;
