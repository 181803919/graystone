
#include "fly_stl_predefine.h"

#ifdef USE_MYSKIP
	#include "skip/skip_list.h"
#endif

int main()
{
	srand(time(NULL));

	typedef Skip_List<int, Compare_Fun<int>, Equal_Fun<int>> Skip_List_Int;
	Skip_List_Int skip_list_test;
	skip_list_test.Reserve(1000000);
	//int test_array[10] = { 3,4,2,23,24,29,67,1,5,52};
	for (int i = 0; i < 1000000; ++i)
	{
		//int value = (rand() % 10000);
		if (i % 10000 == 0)
		{
			printf("Init %d\n", i);
		}
		
		if (NULL == skip_list_test.InsertNode(i))
		{
			break;
		}
	}

	printf("free node:%llu\n", skip_list_test.mem_list_.size());
	for (int i = 0; i < 1000000; ++i)
	{
		//int value = (rand() % 10000);
		if (i % 10000 == 0)
		{
			printf("del %d\n", i);
			printf("free node:%llu\n", skip_list_test.mem_list_.size());
		}

		skip_list_test.DeleteNode(i);
	}

	//typedef std::set<int> Set_Int;
	//Set_Int tmp_set;
	//for (int i = 0; i < 100000; ++i)
	//{
	//	if (i % 10000 == 0)
	//	{
	//		printf("Init %d\n", i);
	//	}
	//	tmp_set.insert(i);
	//}
	skip_list_test.Foreach(&PrintSkipNodeInt);
	printf("free node:%llu\n", skip_list_test.mem_list_.size());

	//skip_list_test.SetDebug(true);
	//skip_list_test.SetPrintFun(&PrintSkipNodeInt);
	//for (int i = 0; i < 1; ++i)
	//{
	//	int value = rand() % 10000 + 10000;
	//	printf("Insert: %d\n", value);
	//	skip_list_test.InsertNode(value);
	//}

	//skip_list_test.Foreach(&PrintSkipNodeInt);
	printf("Init Over\n");

	typedef Skip_List_Node<int, Compare_Fun<int>, Equal_Fun<int>> Skip_List_Node_Int;
	int x = 6584;
	//Skip_List_Node_Int* search = skip_list_test.SeachNode(x);
	//skip_list_test.Foreach(&PrintSkipNodeInt);


	int a = getchar();
	return 0;
}