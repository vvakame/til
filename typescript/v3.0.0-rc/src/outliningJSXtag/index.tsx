const Hoge = () => <>hoge</>;
const Fuga = () => <Hoge></Hoge>;
const Piyo = () => <Piyo></Piyo>;

const V = () => {
    return (
        <Hoge>
            <Fuga>
                <Piyo />
            </Fuga>
        </Hoge>
    )
};
