type DeepPartial<T> = { [P in keyof T]?: DeepPartial<T[P]> };
type Writeable<T> = { -readonly [P in keyof T]-?: Writeable<T[P]> };

export type DataConstructor<T> = DeepPartial<Writeable<T>>;
