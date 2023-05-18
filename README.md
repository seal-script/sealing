# sealing
The interpreter(compiler other than C) for SealScript

# example
```haskell
module example (
    fact,
    List(..),
    Category(..),
)


fact : Int -> Int
fact 0 = 0
fact 1 = 1
fact n = fact (n - 1) + fact (n - 2)

double (x : Int) = x * x

List : Type -> Type
enum List a {
    Nil  : List a
    (::) : a -> List a -> List a
}

enum Expr {
    Number { value : Float }
    FuncCall { f : String, args : List Expr }
}

enum Person {
    New { id : Int, name : String }
    OfId Int
}

seal Monoid a {
    zero : a
    (<>) : a -> a -> a
}

seal Eq a {
    (==) : a -> a -> Bool
    x == y = not (x != y)

    (!=) : a -> a -> Bool
    x != y = not (x == y)
}


tom = Person.New {
    id = 0
    name = "Tom"
}

impl Person = tom


showPerson : Person -> String
showPerson p = 
    printf "(Person %d %s)" p.id p.name

f >> g = f (g x)

Category : (Type -> Type -> Type) -> Type
seal Category c {
    id : c a a
    (~) : c a b -> c b c -> c a c
}

impl Category (->) {
    id = \a -> a
    (~) = (>>)
}


Semi : Type -> Type
seal Semi a {
    (<>) : a -> a -> a
}
Monoid : Type -> Type
seal Semi a => Monoid a {
    empty : a
}


(++) : a => List a -> List a -> List a
xs ++ ys = case xs of
    Nil -> ys
    (x :: xs) -> x :: (xs ++ ys)

impl (a : Type) => Semi (List a) {
    xs <> ys = xs ++ ys
}

impl a => Monoid (List a) {
    empty = Nil
}


Functor : (Type -> Type) -> Type
seal Functor f {
    map : (a -> b) -> f a -> f b
}
(<$>) = Functor.map

impl ListFunctor : Functor List {
    map f [] = []
    map f (x :: xs) = (f x) :: map f xs
}


sum : Monoid a => List a -> a
sum [] = empty
sum xs = Ref.run $ do
    ref <- Ref.new empty
    for xs $ \x ->
        Ref.set ref (<> x)
    Ref.get ref


main : IO ()
main = print "Hello, world!"

clear : Ref Person -> Ref Person
clear person = 
    Ref.set person.id (const 0)

Name : Type
Name = String

Collection : Type -> Type
Collection a = List a

Vec : Type -> Int -> Type
seal Vec a n {
    Nil : Vec a 0
    (:+) : a -> Vec a n -> Vec a (n + 1)
}


seal Show a {
    show : a -> String
}

Showable = Show a => a
showIt : Showable -> String
showIt s = show s


Lift : Type -> Type
Lift a = case a of
    Int   -> Long
    Float -> Double

```







